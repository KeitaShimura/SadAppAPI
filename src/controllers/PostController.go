package controllers

import (
	"SadApp/src/database"
	"SadApp/src/middlewares"
	"SadApp/src/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func Posts(c *fiber.Ctx) error {
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

    // Fetch posts from the current user and the users who follow them
    var posts []models.Post
    result := database.DB.Where("user_id IN ?", followerIds).
        Preload("User").
        Limit(pageSize).
        Offset((page - 1) * pageSize).
        Find(&posts)

    if result.Error != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Cannot retrieve posts",
        })
    }

    // Return the list of posts as JSON
    return c.JSON(posts)
}

// Helper function to get pagination parameters
func getPaginationParameters(c *fiber.Ctx) (int, int) {
    page, _ := strconv.Atoi(c.Query("page", "1"))
    pageSize, _ := strconv.Atoi(c.Query("pageSize", "50"))

    if page < 1 {
        page = 1
    }
    if pageSize < 1 || pageSize > 100 {
        pageSize = 50
    }

    return page, pageSize
}


func UserPosts(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	var posts []models.Post
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
	database.DB.Model(&models.Post{}).Where("user_id = ?", userID).Count(&total)

	result := database.DB.Where("user_id = ?", userID).
		Preload("User").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&posts)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot retrieve posts for the user",
		})
	}

	return c.JSON(posts)
}

func CreatePost(c *fiber.Ctx) error {
	// First, get the user ID from the JWT token
	userId, err := middlewares.GetUserId(c)
	if err != nil {
		// If there's an error retrieving the user ID, return an error
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Initialize a new Post struct
	var post models.Post

	// Parse the request body into the post struct
	if err := c.BodyParser(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad request",
		})
	}

	// Assign the retrieved user ID to the post
	post.UserID = userId // Assuming your Post model has a UserId field

	// Create the post in the database
	result := database.DB.Create(&post)
	if result.Error != nil {
		// If there's an error during the creation, return the error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot create post",
		})
	}

	// Return the created post as JSON
	return c.JSON(post)
}

func GetPost(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	post := models.Post{
		Id: uint(id),
	}

	database.DB.Find(&post)

	return c.JSON(post)
}

func UpdatePost(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	post := models.Post{
		Id: uint(id),
	}

	if err := c.BodyParser(&post); err != nil {
		return err
	}

	database.DB.Model(&post).Updates(post)

	return c.JSON(post)
}

func DeletePost(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	post := models.Post{
		Id: uint(id),
	}

	database.DB.Delete(&post)

	return nil
}
