package stock

import (
	"context"
	"encoding/csv"
	"io"
	"strconv"
	"strings"

	"Veritasbackend/internal/domain/repositories"
)

type UploadProductsUseCase struct {
	productRepo repositories.ProductRepository
}

func NewUploadProductsUseCase(productRepo repositories.ProductRepository) *UploadProductsUseCase {
	return &UploadProductsUseCase{
		productRepo: productRepo,
	}
}

type UploadResult struct {
	Imported int      `json:"imported"`
	Errors   []string `json:"errors"`
}

func (uc *UploadProductsUseCase) Execute(ctx context.Context, tenantID int, reader io.Reader) (*UploadResult, error) {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = ','
	csvReader.Comment = '#'

	// Leer header
	header, err := csvReader.Read()
	if err != nil {
		return nil, err
	}

	// Validar header esperado: name,description,price,stock,sku
	if len(header) < 3 {
		return &UploadResult{Errors: []string{"Invalid CSV format"}}, nil
	}

	imported := 0
	errors := []string{}

	// Leer registros
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			errors = append(errors, err.Error())
			continue
		}

		if len(record) < 3 {
			errors = append(errors, "Invalid record format")
			continue
		}

		// Parsear campos
		name := strings.TrimSpace(record[0])
		description := ""
		if len(record) > 1 {
			description = strings.TrimSpace(record[1])
		}

		price, err := strconv.ParseFloat(strings.TrimSpace(record[2]), 64)
		if err != nil {
			errors = append(errors, "Invalid price: "+record[2])
			continue
		}

		stock := 0
		if len(record) > 3 {
			stock, err = strconv.Atoi(strings.TrimSpace(record[3]))
			if err != nil {
				stock = 0
			}
		}

		sku := ""
		if len(record) > 4 {
			sku = strings.TrimSpace(record[4])
		}

		// Crear producto
		_, err = uc.productRepo.Create(ctx, tenantID, name, description, sku, price, stock)
		if err != nil {
			errors = append(errors, "Failed to create product: "+name)
			continue
		}

		imported++
	}

	return &UploadResult{
		Imported: imported,
		Errors:   errors,
	}, nil
}

