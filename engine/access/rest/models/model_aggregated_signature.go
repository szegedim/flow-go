/*
 * Access API
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package models

type AggregatedSignature struct {
	VerifierSignatures []string `json:"verifier_signatures"`

	SignerIds []string `json:"signer_ids"`
}
