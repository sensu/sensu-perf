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

	_ "github.com/go-sql-driver/mysql"
)

var (
	concurrency   = flag.Int("c", 50, "max number of read/write goroutines")
	conns         = flag.Int("conns", 24, "max number of connections")
	numEvents     = flag.Int("n", 1000, "total number of events")
	myURL         = flag.String("d", "root@tcp(127.0.0.1:3306)/sensu", "mysql url")
	atomicCounter int64
)

const ddl = `CREATE TABLE IF NOT EXISTS events (
    id              serial       PRIMARY KEY,
    sensu_namespace varchar(256) NOT NULL,
    sensu_entity    varchar(256) NOT NULL,
    sensu_check     varchar(256) NOT NULL,
    serialized      mediumblob   NOT NULL,
    UNIQUE KEY( sensu_namespace, sensu_entity, sensu_check )
) ENGINE = INNODB;`

const selectEventQuery = `SELECT serialized FROM events
    WHERE sensu_namespace = ? AND sensu_entity = ? AND sensu_check = ?`

const insertEventQuery = `INSERT INTO events (
        sensu_namespace,
        sensu_entity,
        sensu_check,
        serialized
    )
    VALUES (
        ?,
        ?,
        ?,
        ?
    )`

const updateEventQuery = `UPDATE events
    SET serialized = ?
    WHERE sensu_namespace = ? AND sensu_entity = ? AND sensu_check = ?`

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

	db, err := sql.Open("mysql", *myURL)
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

func withTx(ctx context.Context, db *sql.DB, fun func(*sql.Tx) error) (err error) {
	tx, bErr := db.BeginTx(ctx, nil)
	if bErr != nil {
		return bErr
	}

	defer func() {
		if err == nil {
			err = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
	}()

	return fun(tx)
}

func bencher(ctx context.Context, db *sql.DB, sel *sql.Stmt, ins *sql.Stmt, up *sql.Stmt, index int, eventNames []string, serialized []byte) {
	j := index
	for ctx.Err() == nil {
		eventName := eventNames[j]
		j = (j + index) % len(eventNames)

		err := withTx(ctx, db, func(tx *sql.Tx) error {
			var result []byte
			sel := tx.Stmt(sel)
			err := sel.QueryRowContext(ctx, "default", eventName, eventName).Scan(&result)

			switch err {
			case sql.ErrNoRows:
				_, err := ins.ExecContext(ctx, "default", eventName, eventName, serialized)
				if err != nil && ctx.Err() == nil {
					return err
				}
			case nil:
				_, err := up.ExecContext(ctx, "default", eventName, eventName, serialized)
				if err != nil && ctx.Err() == nil {
					return err
				}
			default:
				return err
			}

			return nil
		})

		if err == nil {
			countEvent()
		} else {
			log.Printf("%s", err)
		}
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

	sel, err := db.PrepareContext(ctx, selectEventQuery)
	if err != nil && err != context.Canceled {
		return fmt.Errorf("error preparing select statement: %s", err)
	}
	defer sel.Close()

	ins, err := db.PrepareContext(ctx, insertEventQuery)
	if err != nil && err != context.Canceled {
		return fmt.Errorf("error preparing insert statement: %s", err)
	}
	defer ins.Close()

	up, err := db.PrepareContext(ctx, updateEventQuery)
	if err != nil && err != context.Canceled {
		return fmt.Errorf("error preparing update statement: %s", err)
	}
	defer up.Close()

	for i := 0; i < *concurrency; i++ {
		go bencher(ctx, db, sel, ins, up, i, eventNames, serialized)
	}

	<-ctx.Done()

	return nil
}
