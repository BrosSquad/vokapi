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

type (
	ErrorReponse struct {
		Message string `json:"message,omitempty"`
		Name    string `json:"string,omitempty"`
	}
	Response struct {
		Vocative string `json:"vokativ"`
	}
)


func Register(di *container.Container, router fiber.Router) {
	router.Get("/:name", Handler(di.GetBadgerDB()))
}

func Handler(db *badger.DB) fiber.Handler {
	caser := cases.Title(language.SerbianLatin)

	return func(c *fiber.Ctx) error {
		name := c.Params("name", "")
		decodedName, err := url.QueryUnescape(name)
		if err != nil {
			return c.
				Status(fiber.StatusUnprocessableEntity).
				JSON(ErrorReponse{
					Message: "ime nije validno",
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
				return c.
					Status(fiber.StatusNotFound).
					JSON(ErrorReponse{
						Message: "ime nije pronadjeno",
						Name:    decodedName,
					})
			}

			return c.
				Status(fiber.StatusInternalServerError).
				JSON(ErrorReponse{
					Message: "Ups, nešto je pošlo po zlu.",
				})
		}

		return c.
			Status(fiber.StatusOK).
			JSON(Response{
				Vocative: utils.UnsafeString(caser.Bytes(value)),
			})
	}
}
