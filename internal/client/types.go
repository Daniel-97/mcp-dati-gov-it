package client

// CKANResponse è l'envelope comune a tutte le risposte CKAN.
type CKANResponse[T any] struct {
	Success bool       `json:"success"`
	Result  T          `json:"result"`
	Error   *CKANError `json:"error,omitempty"`
}

// CKANError è il payload di errore restituito da CKAN.
type CKANError struct {
	Message string `json:"message"`
	Type    string `json:"__type"`
}

// SearchResult è il payload di package_search.
type SearchResult struct {
	Count   int       `json:"count"`
	Results []Dataset `json:"results"`
}

// Dataset rappresenta un package CKAN (dataset).
type Dataset struct {
	ID               string        `json:"id"`
	Name             string        `json:"name"`
	Title            string        `json:"title"`
	Notes            string        `json:"notes"`
	Tags             []Tag         `json:"tags"`
	Organization     *Organization `json:"organization"`
	Resources        []Resource    `json:"resources"`
	MetadataModified string        `json:"metadata_modified"`
	LicenseTitle     string        `json:"license_title"`
}

// Tag è un tag CKAN.
type Tag struct {
	Name string `json:"name"`
}

// Organization è un'organizzazione CKAN (PA).
type Organization struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Resource è un file allegato a un dataset.
type Resource struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Format string `json:"format"`
	URL    string `json:"url"`
	Size   int64  `json:"size"`
}
