package controllers

import (
	"server/database"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func HandleGetSensorsWithPagination(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))

	sensors, totalRecords, err := database.GetSensorDataWithPagination(page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	totalPages := (totalRecords + int64(pageSize) - 1) / int64(pageSize)
	hasNext := int64(page*pageSize) < totalRecords
	hasPrev := page > 1
	return c.JSON(fiber.Map{
		"page":         page,
		"limit":        pageSize,
		"totalPages":   totalPages,
		"hasNext":      hasNext,
		"hasPrev":      hasPrev,
		"totalRecords": totalRecords,
		"data":         sensors,
	})
}
