package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/google/uuid"
	sensu "github.com/sensu/sensu-go/api/core/v2"

	_ "github.com/lib/pq"
)

var (
	concurrency = flag.Int("c", 50, "max number of read/write goroutines")
	conns       = flag.Int("conns", 24, "max number of connections")
	numEvents   = flag.Int("n", 1000, "total number of events")
	pgURL       = flag.String("d", "host=/run/postgresql sslmode=disable", "postgresql url")

	atomicCounter int64
)

const ddl = `CREATE TABLE IF NOT EXISTS events (
    id              serial     PRIMARY KEY,
    sensu_namespace text       NOT NULL,
    sensu_check     text       NOT NULL,
    sensu_entity    text       NOT NULL,
    history_index   integer    DEFAULT 2,
	history_ts      integer[],
	history_status  integer[],
    serialized      jsonb, 
    UNIQUE ( sensu_namespace, sensu_check, sensu_entity )
);`

const writeEventQuery = `INSERT INTO events (
        sensu_namespace,
        sensu_entity,
        sensu_check,
        status,
        last_ok,
        occurrences,
        occurrences_wm,
        history_status,
        history_ts,
        serialized
    )
    VALUES (
        $1,
        $2,
        $3,
        $4,
        CASE
            WHEN $4 = 0 THEN $5
            ELSE 0
        END,
        1,
        1,
        ARRAY[$4]::integer[],
        ARRAY[$5]::integer[],
        $6
    )
    ON CONFLICT (
        sensu_namespace,
        sensu_entity,
        sensu_check
    ) DO UPDATE SET (
        status,
        last_ok,
        occurrences,
        occurrences_wm,
        history_status[events.history_index],
        history_ts[events.history_index],
        history_index,
        serialized,
	previous_serialized
    ) = (
        $4,
        CASE
            WHEN $4 = 0 THEN $5
            ELSE events.last_ok
        END,
        CASE
            WHEN events.status = $4 THEN events.occurrences + 1
            ELSE 1
        END,
        CASE
            WHEN events.status = 0 AND $4 != 0 THEN 1
            WHEN events.status != $4 THEN events.occurrences_wm
            WHEN events.occurrences < events.occurrences_wm THEN events.occurrences_wm
            ELSE events.occurrences_wm + 1
        END,
        $4,
        $5,
        (events.history_index % 20) + 1,
        $6,
	events.serialized
    )
    RETURNING history_ts, history_status, last_ok, occurrences, occurrences_wm, previous_serialized;`

func initTest(db *sql.DB) error {
	_, err := db.Exec(ddl)
	if err != nil {
		return fmt.Errorf("couldn't initialize test: %s", err)
	}
	return nil
}

func deleteTest(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM events;")
	if err != nil {
		return fmt.Errorf("error cleaning up test: %s", err)
	}
	return nil
}

func main() {
	flag.Parse()

	db, err := sql.Open("postgres", *pgURL)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(*conns)

	if err := initTest(db); err != nil {
		log.Fatal(err)
	}

	defer deleteTest(db)

	go eventCounter()

	if err := doBench(db); err != nil {
		log.Fatal(err)
	}
}

func countEvent() {
	atomic.AddInt64(&atomicCounter, 1)
}

func eventCounter() {
	ticker := time.NewTicker(time.Second)
	for range ticker.C {
		counted := atomic.SwapInt64(&atomicCounter, 0)
		fmt.Printf("%d events/sec\n", counted)
	}
}

func bencher(ctx context.Context, write *sql.Stmt, index int, eventNames []string, serialized []byte) {
	j := index
	for ctx.Err() == nil {
		eventName := eventNames[j]
		j = (j + index) % len(eventNames)
		_, err := write.ExecContext(ctx, "default", eventName, eventName, 1, 1, serialized)
		if err != nil && ctx.Err() == nil {
			log.Printf("error executing write: %s", err)
			continue
		}
		countEvent()
	}
}

const cannedResponseText = `
                         .'loo:,
                        ,KNMMWNWX
                  ..    ,000OkxkW'
                 ,o,.   .O0KOOOk0:
                 dkl.    :OO0kkkk;
                 :ko     .lOOOOkx
              'OXX0:       'xWMdkd;o,.
              cMMMM; .,lkXN;oWNOk,;MMMMWXx:
              oMMMM:KNWMMMMl'.cl .NMMMMMMMMX
              NMMMWkMMMMMMMMWxxKONMMMMMMMMMM
             oMMMMMMMMMMMMMMMNW0NMMMMMMMMMMM.
             KMMMMMMMMMMMMMMWMWWMMMMMMMMMMMMN
             oKXXKKKKXMMMMMMMMMWMMMMMMMMMMMMM.
                     'MMMMMMMMMMMMMMMMMMMMMMMk
                     .MMMMMMMMWMMMMMMMMMMMMMMM
                      WMMMMMMMMMMWNMMMMMMMMMMN
                      WMMMMMMMMMWX0kO0WMMMMMMO
                     .MMMMMMMMMMMNX0kkWMMMMWO'
                     ;MMMMMMMMMMMMWXNNMMMMW.
`

func doBench(db *sql.DB) error {
	event := sensu.FixtureEvent("entity", "check")
	event.Check.Output = cannedResponseText
	serialized, err := proto.Marshal(event)
	if err != nil {
		return err
	}

	eventNames := make([]string, *numEvents)

	for i := 0; i < *numEvents; i++ {
		eventNames[i] = uuid.New().String()
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-sigc
		cancel()
	}()

	write, err := db.PrepareContext(ctx, writeEventQuery)
	if err != nil && err != context.Canceled {
		return fmt.Errorf("error preparing write statement: %s", err)
	}
	defer write.Close()

	for i := 0; i < *concurrency; i++ {
		go bencher(ctx, write, i, eventNames, serialized)
	}

	<-ctx.Done()

	return nil
}
