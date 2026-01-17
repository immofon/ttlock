package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/immofon/ttlock"
	"github.com/urfave/cli/v2"
)

func printJSON(v interface{}) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

func parseDate(s string) (int64, error) {
	// Format: "20230214-14" -> YYYYMMDD-HH
	layout := "20060102-15"
	t, err := time.ParseInLocation(layout, s, time.Local)
	if err != nil {
		return 0, err
	}
	return t.UnixMilli(), nil
}

var lockCmd = &cli.Command{
	Name:  "lock",
	Usage: "Get lock details",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:     "id",
			Required: true,
			Usage:    "Lock ID",
		},
	},
	Action: func(c *cli.Context) error {
		lockID := c.Int("id")
		detail, err := client.GetLockDetail(lockID)
		if err != nil {
			return err
		}
		return printJSON(detail)
	},
}

var listLockCmd = &cli.Command{
	Name:  "list-lock",
	Usage: "List locks",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  "n",
			Usage: "Page number",
			Value: 1,
		},
		&cli.IntFlag{
			Name:  "s",
			Usage: "Page size",
			Value: 20,
		},
		&cli.StringFlag{
			Name:  "a",
			Usage: "Lock alias",
		},
		&cli.IntFlag{
			Name:  "g",
			Usage: "Group ID",
		},
	},
	Action: func(c *cli.Context) error {
		pageNo := c.Int("n")
		pageSize := c.Int("s")
		lockAlias := c.String("a")
		groupID := c.Int("g")

		list, err := client.GetLockList(pageNo, pageSize, lockAlias, groupID)
		if err != nil {
			return err
		}
		return printJSON(list)
	},
}

var listPasscodeCmd = &cli.Command{
	Name:  "list-passcode",
	Usage: "List passcodes",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:     "id",
			Required: true,
			Usage:    "Lock ID",
		},
		&cli.IntFlag{
			Name:  "n",
			Usage: "Page number",
			Value: 1,
		},
		&cli.IntFlag{
			Name:  "s",
			Usage: "Page size",
			Value: 20,
		},
		&cli.IntFlag{
			Name:  "o",
			Usage: "Order by (1:desc, 2:asc)",
			Value: 1,
		},
		&cli.StringFlag{
			Name:    "search",
			Aliases: []string{"q"},
			Usage:   "Search string",
		},
	},
	Action: func(c *cli.Context) error {
		lockID := c.Int("id")
		pageNo := c.Int("n")
		pageSize := c.Int("s")
		orderBy := c.Int("o")
		searchStr := c.String("search")

		list, err := client.GetPasscodeList(lockID, pageNo, pageSize, orderBy, searchStr)
		if err != nil {
			return err
		}
		return printJSON(list)
	},
}

var genPassCmd = &cli.Command{
	Name:  "genpass",
	Usage: "Generate random passcode",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:     "id",
			Required: true,
			Usage:    "Lock ID",
		},
		&cli.IntFlag{
			Name:     "t",
			Required: true,
			Usage:    "Passcode type",
		},
		&cli.StringFlag{
			Name:  "n",
			Usage: "Passcode name",
		},
		&cli.StringFlag{
			Name:     "s",
			Required: true,
			Usage:    "Start date (YYYYMMDD-HH)",
		},
		&cli.StringFlag{
			Name:     "e",
			Required: true,
			Usage:    "End date (YYYYMMDD-HH)",
		},
	},
	Action: func(c *cli.Context) error {
		lockID := c.Int("id")
		pwdType := ttlock.PasscodeType(c.Int("t"))
		pwdName := c.String("n")
		startDateStr := c.String("s")
		endDateStr := c.String("e")

		startDate, err := parseDate(startDateStr)
		if err != nil {
			return fmt.Errorf("invalid start date: %w", err)
		}
		endDate, err := parseDate(endDateStr)
		if err != nil {
			return fmt.Errorf("invalid end date: %w", err)
		}

		resp, err := client.GetRandomPasscode(lockID, pwdType, pwdName, startDate, endDate)
		if err != nil {
			return err
		}
		return printJSON(resp)
	},
}

var sendKeyCmd = &cli.Command{
	Name:  "sendkey",
	Usage: "Send eKey",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:     "id",
			Required: true,
			Usage:    "Lock ID",
		},
		&cli.StringFlag{
			Name:     "to",
			Required: true,
			Usage:    "Receiver username",
		},
		&cli.StringFlag{
			Name:  "n",
			Usage: "Key name",
		},
		&cli.StringFlag{
			Name:     "s",
			Required: true,
			Usage:    "Start date (YYYYMMDD-HH)",
		},
		&cli.StringFlag{
			Name:     "e",
			Required: true,
			Usage:    "End date (YYYYMMDD-HH)",
		},
	},
	Action: func(c *cli.Context) error {
		lockID := c.Int("id")
		receiverUsername := c.String("to")
		keyName := c.String("n")
		startDateStr := c.String("s")
		endDateStr := c.String("e")

		startDate, err := parseDate(startDateStr)
		if err != nil {
			return fmt.Errorf("invalid start date: %w", err)
		}
		endDate, err := parseDate(endDateStr)
		if err != nil {
			return fmt.Errorf("invalid end date: %w", err)
		}

		// Using nil for options for now as not specified in CLI flags
		resp, err := client.SendKey(lockID, receiverUsername, keyName, startDate, endDate, nil)
		if err != nil {
			return err
		}
		return printJSON(resp)
	},
}
