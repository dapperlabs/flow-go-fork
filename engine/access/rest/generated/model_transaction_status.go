/*
 * Access API
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package generated

type TransactionStatus string

// List of TransactionStatus
const (
	PENDING   TransactionStatus = "Pending"
	FINALIZED TransactionStatus = "Finalized"
	EXECUTED  TransactionStatus = "Executed"
	SEALED    TransactionStatus = "Sealed"
	EXPIRED   TransactionStatus = "Expired"
)