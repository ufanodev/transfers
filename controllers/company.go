package controllers

import (
	"net/http"
	"transfers/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CompanyController struct {
	DB *gorm.DB
}

// GET /api/v1/companies
// Recupera todas las empresas ordenadas alfabéticamente
func (ctrl *CompanyController) GetAll(c *gin.Context) {
	var items []models.Company
	if err := ctrl.DB.Order("name asc").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al recuperar listado de empresas"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GET /api/v1/companies/:id
// Obtiene el detalle de una empresa para cargar en el formulario de edición
func (ctrl *CompanyController) Get(c *gin.Context) {
	var item models.Company
	id := c.Param("id")

	if err := ctrl.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Empresa no encontrada"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// POST /api/v1/companies
// Crea una nueva empresa vinculándola a un usuario responsable
func (ctrl *CompanyController) POST(c *gin.Context) {
	var item models.Company
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos del formulario inválidos"})
		return
	}

	// VALIDACIÓN: Verificar que se ha enviado un user_id
	if item.UserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Debe seleccionar un usuario responsable para esta empresa"})
		return
	}

	// Transacción para asegurar la integridad
	tx := ctrl.DB.Begin()

	if err := tx.Create(&item).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error de base de datos al guardar la empresa"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusCreated, item)
}

// PUT /api/v1/companies/:id
// Actualiza los datos de una empresa existente
func (ctrl *CompanyController) PUT(c *gin.Context) {
	var item models.Company
	id := c.Param("id")

	// 1. Verificar existencia del registro
	if err := ctrl.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Empresa no encontrada para actualizar"})
		return
	}

	// 2. Vincular nuevos datos del JSON
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de actualización inválidos"})
		return
	}

	// 3. Guardar cambios (GORM Save actualiza todos los campos del modelo)
	if err := ctrl.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar los cambios en la base de datos"})
		return
	}

	c.JSON(http.StatusOK, item)
}

// DELETE /api/v1/companies/:id
// Elimina una empresa (Soporta Soft Delete si el modelo tiene DeletedAt)
func (ctrl *CompanyController) DELETE(c *gin.Context) {
	id := c.Param("id")

	// Verificamos si la empresa existe antes de intentar borrar
	var item models.Company
	if err := ctrl.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "La empresa no existe"})
		return
	}

	if err := ctrl.DB.Delete(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo eliminar la empresa"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Empresa eliminada correctamente"})
}
