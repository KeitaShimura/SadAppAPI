package controllers

import (
	"SadApp/src/database"
	"SadApp/src/middlewares"
	"SadApp/src/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func Events(c *fiber.Ctx) error {
	currentUserId, err := middlewares.GetUserId(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error retrieving user ID",
		})
	}

	// Fetch the IDs of users that are following the current user
	var followers []models.Follow
	database.DB.Where("following_id = ?", currentUserId).Find(&followers)

	// Extract user IDs from the followers
	var followerIds []uint
	for _, follower := range followers {
		followerIds = append(followerIds, follower.FollowerId)
	}

	// Include the current user's ID in the list
	followerIds = append(followerIds, currentUserId)

	// Get pagination parameters
	page, pageSize := getPaginationParameters(c)

	// Fetch events from the current user and the users who follow them
	var events []models.Event
	result := database.DB.Where("user_id IN ?", followerIds).
		Preload("User").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&events)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot retrieve events",
		})
	}

	// Return the list of events as JSON
	return c.JSON(events)
}

func UserEvents(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var events []models.Event
	// ページ番号を取得
	page := 1
	pageSize := 50

	// クエリから 'page' を取得
	if p, err := strconv.Atoi(c.Query("page", "1")); err == nil && p > 0 {
		page = p
	}

	// クエリから 'pageSize' を取得
	if ps, err := strconv.Atoi(c.Query("pageSize", "50")); err == nil && ps > 0 {
		pageSize = ps
	}

	var total int64
	database.DB.Model(&models.Event{}).Where("user_id = ?", userID).Count(&total)

	result := database.DB.Where("user_id = ?", userID).
		Preload("User").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&events)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot retrieve events for the user",
		})
	}

	return c.JSON(events)
}

func GetEvent(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	event := models.Event{
		Id: uint(id),
	}

	database.DB.Find(&event)

	return c.JSON(event)
}

func CreateEvent(c *fiber.Ctx) error {
	// First, get the user ID from the JWT token
	userId, err := middlewares.GetUserId(c)
	if err != nil {
		// If there's an error retrieving the user ID, return an error
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Initialize a new event struct
	var event models.Event

	// Parse the request body into the event struct
	if err := c.BodyParser(&event); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad request",
		})
	}

	// Assign the retrieved user ID to the event
	event.UserId = userId // Assuming your event model has a UserId field

	// Create the event in the database
	result := database.DB.Create(&event)
	if result.Error != nil {
		// If there's an error during the creation, return the error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot create event",
		})
	}

	database.DB.Preload("User").Find(&event, event.Id)

	// Return the created event as JSON
	return c.JSON(event)
}

func UpdateEvent(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	event := models.Event{
		Id: uint(id),
	}

	if err := c.BodyParser(&event); err != nil {
		return err
	}

	database.DB.Model(&event).Updates(event)

	return c.JSON(event)
}

func UserLikedEvents(c *fiber.Ctx) error {
	// Retrieve the user ID (adjust this part based on how you manage user identification)
	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Get pagination parameters (reuse your existing function or define one)
	page, pageSize := getPaginationParameters(c)

	// Find IDs of events liked by the user
	var eventLikes []models.EventLike
	database.DB.Where("user_id = ?", userID).Find(&eventLikes)

	// Extract event IDs
	var eventIds []uint
	for _, eventLike := range eventLikes {
		eventIds = append(eventIds, eventLike.EventId)
	}

	// Fetch the events based on the event IDs
	var events []models.Event
	result := database.DB.Where("id IN ?", eventIds).
		Preload("User").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&events)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot retrieve liked events",
		})
	}

	// Return the list of liked events as JSON
	return c.JSON(events)
}

func DeleteEvent(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	event := models.Event{
		Id: uint(id),
	}

	database.DB.Delete(&event)

	return nil
}
