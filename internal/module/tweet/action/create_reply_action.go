package action

import (
	"errors"
	"strconv"

	"github.com/HotPotatoC/twitter-clone/internal/module"
	"github.com/HotPotatoC/twitter-clone/internal/module/tweet/entity"
	"github.com/HotPotatoC/twitter-clone/internal/module/tweet/service"
	"github.com/gofiber/fiber/v2"
)

type createReplyAction struct {
	service service.CreateReplyService
}

func NewCreateReplyAction(service service.CreateReplyService) module.Action {
	return createReplyAction{service: service}
}

func (a createReplyAction) Execute(c *fiber.Ctx) error {
	var input service.CreateReplyInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if errors := input.Validate(); errors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	tweetID, err := strconv.ParseInt(c.Params("tweetID"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "There was a problem on our side",
		})
	}
	userID := c.Locals("userID").(float64)

	err = a.service.Execute(input, int64(userID), tweetID)
	if err != nil {
		if errors.Is(err, entity.ErrTweetDoesNotExist) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Tweet not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "There was a problem on our side",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Successfully posted a new reply",
	})
}
