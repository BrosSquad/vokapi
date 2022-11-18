package vokativ

import (
	"errors"
	"net/url"

	"github.com/dgraph-io/badger/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/BrosSquad/vokapi/container"
)

func Register(di *container.Container, router fiber.Router) {
	router.Get("/:name", handler(di.GetBadgerDB()))
}

func handler(db *badger.DB) fiber.Handler {
	caser := cases.Title(language.SerbianLatin)

	return func(c *fiber.Ctx) error {
		name := c.Params("name", "")

		decodedName, err := url.QueryUnescape(name)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"poruka": "ime nije validno",
			})
		}

		var value []byte

		err = db.View(func(txn *badger.Txn) error {
			item, err := txn.Get(utils.UnsafeBytes(caser.String(decodedName)))
			if err != nil {
				return err
			}

			value, err = item.ValueCopy(nil)
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"poruka": "Ime nije pronađeno.",
					"ime":    decodedName,
				})
			}

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"poruka": "Ups, nešto je pošlo po zlu.",
			})
		}
		// golang.org/x/text/cases

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"vokativ": string(caser.Bytes(value)),
		})
	}
}
