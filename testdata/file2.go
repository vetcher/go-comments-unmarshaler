package testdata

// FetchEmails is for fetching emails.
//
// This file is placed in second file.
func (Fetcher) FetchEmails() error { return nil }

/*
fetchUsers is private function.
*/
func (*Fetcher) fetchUsers() error { return nil }
