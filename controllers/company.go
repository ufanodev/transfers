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

// GET /api/v1/companies - Listar todas las empresas
func (ctrl *CompanyController) GetAll(c *gin.Context) {
	var items []models.Company
	// Preload("User") para ver los datos de cuenta asociados
	if err := ctrl.DB.Preload("User").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al recuperar empresas"})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GET /api/v1/companies/:id - Obtener una empresa específica
func (ctrl *CompanyController) Get(c *gin.Context) {
	var item models.Company
	if err := ctrl.DB.Preload("User").First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Empresa no encontrada"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// POST /api/v1/companies - Crear nueva empresa
func (ctrl *CompanyController) POST(c *gin.Context) {
	var item models.Company
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear la empresa"})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// PUT /api/v1/companies/:id - Actualizar datos de la empresa
func (ctrl *CompanyController) PUT(c *gin.Context) {
	var item models.Company
	if err := ctrl.DB.First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Empresa no encontrada"})
		return
	}

	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctrl.DB.Save(&item)
	c.JSON(http.StatusOK, item)
}

// DELETE /api/v1/companies/:id - Eliminar empresa
func (ctrl *CompanyController) DELETE(c *gin.Context) {
	// Nota: Esto borrará la empresa, pero el User asociado permanece
	// a menos que configures ON DELETE CASCADE en la DB o lógica extra aquí.
	if err := ctrl.DB.Delete(&models.Company{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar la empresa"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Empresa eliminada correctamente"})
}
