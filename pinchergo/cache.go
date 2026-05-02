package pinchergo

import (
	"maps"
	"sync"
	"time"
)

type cacheEntry[T any] struct {
	Data        *T        `json:"data"`
	EndpointURL string    `json:"endpoint_url"` // endpoint via which this entry was acquired
	CreatedAt   time.Time `json:"created_at"`   // time at which this entry was cached
}

type subcache[T any] map[string]*cacheEntry[T]

// budgetCache is a subcache that stores its own relevant
// budget resources such as accounts, payees, groups, etc.
type budgetCache struct {
	AccountCache    subcache[Account]           `json:"account_cache"`
	PayeeCache      subcache[Payee]             `json:"payee_cache"`
	GroupCache      subcache[Group]             `json:"group_cache"`
	CategoryCache   subcache[Category]          `json:"category_cache"`
	TxnCache        subcache[Transaction]       `json:"transaction_cache"`
	TxnDetailsCache subcache[TransactionDetail] `json:"transaction_details_cache"`
	BudgetEntry     cacheEntry[Budget]          `json:"budget_entry"`
}

// GET
func (c *budgetCache) Budget() *Budget {
	return c.BudgetEntry.Data
}

func (c *budgetCache) Account(aID string) *Account {
	return c.AccountCache[aID].Data
}

func (c *budgetCache) Payee(pID string) *Payee {
	return c.PayeeCache[pID].Data
}

func (c *budgetCache) Group(gID string) *Group {
	return c.GroupCache[gID].Data
}

func (c *budgetCache) Category(cID string) *Category {
	return c.CategoryCache[cID].Data
}

func (c *budgetCache) Transaction(tID string) *Transaction {
	return c.TxnCache[tID].Data
}

func (c *budgetCache) TransactionDetail(tID string) *TransactionDetail {
	return c.TxnDetailsCache[tID].Data
}

type Cache struct {
	mu *sync.Mutex

	// map of budget IDs to a subcache of budget resources
	Entries map[string]*budgetCache `json:"cached_entries"`

	// whether to update cache entries from relevant API calls
	// related to singleton resources
	trackAPICalls bool

	// whether to update cache entries from relevant API calls
	// related to resource collections
	trackBulkAPICalls bool

	interval time.Duration
}

// ----------------------
//    GETTER FUNCTIONS
// ----------------------

func (c *Cache) Budget(bID string) *Budget {
	c.mu.Lock()
	defer c.mu.Unlock()

	b, ok := c.Entries[bID]
	if !ok {
		return nil
	}

	return b.Budget()
}

func (c *Cache) Budgets(urlQuery string) []*Budget {
	c.mu.Lock()
	defer c.mu.Unlock()

	var budgets []*Budget
	for _, b := range c.Entries {
		if b.BudgetEntry.EndpointURL != EndpointBudgets()+urlQuery {
			continue
		}
		budgets = append(budgets, b.BudgetEntry.Data)
	}

	return budgets
}

func (c *Cache) Account(bID, aID string) *Account {
	c.mu.Lock()
	defer c.mu.Unlock()

	b, ok := c.Entries[bID]
	if !ok {
		return nil
	}

	return b.Account(aID)
}

func (c *Cache) Accounts(bID, urlQuery string) []*Account {
	c.mu.Lock()
	defer c.mu.Unlock()

	b, ok := c.Entries[bID]
	if !ok {
		return nil
	}

	accounts := make([]*Account, 0, len(b.AccountCache))
	for _, entry := range b.AccountCache {
		if entry.EndpointURL != EndpointBudgetAccounts(bID)+urlQuery {
			continue
		}
		accounts = append(accounts, entry.Data)
	}

	return accounts
}

func (c *Cache) Payee(bID, pID string) *Payee {
	c.mu.Lock()
	defer c.mu.Unlock()

	b, ok := c.Entries[bID]
	if !ok {
		return nil
	}

	return b.Payee(pID)
}

func (c *Cache) Payees(bID, urlQuery string) []*Payee {
	c.mu.Lock()
	defer c.mu.Unlock()

	b, ok := c.Entries[bID]
	if !ok {
		return nil
	}

	payees := make([]*Payee, 0, len(b.PayeeCache))
	for _, entry := range b.PayeeCache {
		if entry.EndpointURL != EndpointBudgetPayees(bID)+urlQuery {
			continue
		}
		payees = append(payees, entry.Data)
	}

	return payees
}

func (c *Cache) Group(bID, gID string) *Group {
	c.mu.Lock()
	defer c.mu.Unlock()

	b, ok := c.Entries[bID]
	if !ok {
		return nil
	}

	return b.Group(gID)
}

func (c *Cache) Groups(bID, urlQuery string) []*Group {
	c.mu.Lock()
	defer c.mu.Unlock()

	b, ok := c.Entries[bID]
	if !ok {
		return nil
	}

	groups := make([]*Group, 0, len(b.GroupCache))
	for _, entry := range b.GroupCache {
		if entry.EndpointURL != EndpointBudgetGroups(bID)+urlQuery {
			continue
		}
		groups = append(groups, entry.Data)
	}

	return groups
}

func (c *Cache) Category(bID, cID string) *Category {
	c.mu.Lock()
	defer c.mu.Unlock()

	b, ok := c.Entries[bID]
	if !ok {
		return nil
	}

	return b.Category(cID)
}

func (c *Cache) Categories(bID, urlQuery string) []*Category {
	c.mu.Lock()
	defer c.mu.Unlock()

	b, ok := c.Entries[bID]
	if !ok {
		return nil
	}

	categories := make([]*Category, 0, len(b.CategoryCache))
	for _, entry := range b.CategoryCache {
		if entry.EndpointURL != EndpointBudgetCategories(bID)+urlQuery {
			continue
		}
		categories = append(categories, entry.Data)
	}

	return categories
}

func (c *Cache) Transaction(bID, tID string) *Transaction {
	c.mu.Lock()
	defer c.mu.Unlock()

	b, ok := c.Entries[bID]
	if !ok {
		return nil
	}

	return b.Transaction(tID)
}

func (c *Cache) Transactions(bID, urlQuery string) []*Transaction {
	c.mu.Lock()
	defer c.mu.Unlock()

	b, ok := c.Entries[bID]
	if !ok {
		return nil
	}

	txns := make([]*Transaction, 0, len(b.TxnCache))
	for _, entry := range b.TxnCache {
		if entry.EndpointURL != EndpointBudgetTransactions(bID)+urlQuery {
			continue
		}
		txns = append(txns, entry.Data)
	}

	return txns
}

func (c *Cache) TransactionDetails(bID, tID string) *TransactionDetail {
	c.mu.Lock()
	defer c.mu.Unlock()

	b, ok := c.Entries[bID]
	if !ok {
		return nil
	}

	return b.TransactionDetail(tID)
}

func (c *Cache) TransactionsDetails(bID, urlQuery string) []*TransactionDetail {
	c.mu.Lock()
	defer c.mu.Unlock()

	b, ok := c.Entries[bID]
	if !ok {
		return nil
	}

	txns := make([]*TransactionDetail, 0, len(b.TxnDetailsCache))
	for _, entry := range b.TxnDetailsCache {
		if entry.EndpointURL != EndpointBudgetTransactionsDetails(bID)+urlQuery {
			continue
		}
		txns = append(txns, entry.Data)
	}

	return txns
}

// ----------------------
//    SETTER FUNCTIONS
// ----------------------

func (c *Cache) addBudget(dest, bID string, budget *Budget) {
	if !c.trackAPICalls || budget == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	c.Entries[bID] = &budgetCache{
		AccountCache:    map[string]*cacheEntry[Account]{},
		PayeeCache:      map[string]*cacheEntry[Payee]{},
		GroupCache:      map[string]*cacheEntry[Group]{},
		CategoryCache:   map[string]*cacheEntry[Category]{},
		TxnCache:        map[string]*cacheEntry[Transaction]{},
		TxnDetailsCache: map[string]*cacheEntry[TransactionDetail]{},
		BudgetEntry: cacheEntry[Budget]{
			Data:        budget,
			CreatedAt:   time.Now().UTC(),
			EndpointURL: dest,
		},
	}
}

func (c *Cache) addBudgets(dest string, Entries []*Budget) {
	if !c.trackBulkAPICalls || Entries == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, budget := range Entries {
		c.Entries[budget.ID.String()] = &budgetCache{
			AccountCache:    map[string]*cacheEntry[Account]{},
			PayeeCache:      map[string]*cacheEntry[Payee]{},
			GroupCache:      map[string]*cacheEntry[Group]{},
			CategoryCache:   map[string]*cacheEntry[Category]{},
			TxnCache:        map[string]*cacheEntry[Transaction]{},
			TxnDetailsCache: map[string]*cacheEntry[TransactionDetail]{},
			BudgetEntry: cacheEntry[Budget]{
				Data:        budget,
				CreatedAt:   time.Now().UTC(),
				EndpointURL: dest,
			},
		}
	}
}

func (c *Cache) deleteBudget(bID string) {
	if !c.trackAPICalls {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.Entries, bID)
}

func (c *Cache) addAccount(dest, bID string, account *Account) {
	if !c.trackAPICalls || account == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.Entries[bID]; !ok {
		return
	}

	c.Entries[bID].AccountCache[account.ID.String()] = &cacheEntry[Account]{
		Data:        account,
		CreatedAt:   time.Now().UTC(),
		EndpointURL: dest,
	}
}

func (c *Cache) addAccounts(dest, bID string, accounts []*Account) {
	if !c.trackBulkAPICalls || accounts == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.Entries[bID]; !ok {
		return
	}
	for _, a := range accounts {
		if a == nil {
			continue
		}
		c.Entries[bID].AccountCache[a.ID.String()] = &cacheEntry[Account]{
			Data:        a,
			CreatedAt:   time.Now().UTC(),
			EndpointURL: dest,
		}
	}
}

func (c *Cache) deleteAccount(bID, aID string) {
	if !c.trackAPICalls {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.Entries[bID]; !ok {
		return
	}

	delete(c.Entries[bID].AccountCache, aID)
}

func (c *Cache) addPayee(dest, bID string, payee *Payee) {
	if !c.trackAPICalls || payee == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.Entries[bID]; !ok {
		return
	}

	c.Entries[bID].PayeeCache[payee.ID.String()] = &cacheEntry[Payee]{
		Data:        payee,
		CreatedAt:   time.Now().UTC(),
		EndpointURL: dest,
	}
}

func (c *Cache) addPayees(dest, bID string, payees []*Payee) {
	if !c.trackBulkAPICalls || payees == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.Entries[bID]; !ok {
		return
	}
	for _, p := range payees {
		if p == nil {
			continue
		}
		c.Entries[bID].PayeeCache[p.ID.String()] = &cacheEntry[Payee]{
			Data:        p,
			CreatedAt:   time.Now().UTC(),
			EndpointURL: dest,
		}
	}
}

func (c *Cache) deletePayee(bID, pID string) {
	if !c.trackAPICalls {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.Entries[bID]; !ok {
		return
	}

	delete(c.Entries[bID].PayeeCache, pID)
}

func (c *Cache) addGroup(dest, bID string, group *Group) {
	if !c.trackAPICalls || group == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.Entries[bID]; !ok {
		return
	}

	c.Entries[bID].GroupCache[group.ID.String()] = &cacheEntry[Group]{
		Data:        group,
		CreatedAt:   time.Now().UTC(),
		EndpointURL: dest,
	}
}

func (c *Cache) addGroups(dest, bID string, groups []*Group) {
	if !c.trackBulkAPICalls || groups == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.Entries[bID]; !ok {
		return
	}
	for _, g := range groups {
		if g == nil {
			continue
		}
		c.Entries[bID].GroupCache[g.ID.String()] = &cacheEntry[Group]{
			Data:        g,
			CreatedAt:   time.Now().UTC(),
			EndpointURL: dest,
		}
	}
}

func (c *Cache) deleteGroup(bID, gID string) {
	if !c.trackAPICalls {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.Entries[bID]; !ok {
		return
	}

	delete(c.Entries[bID].GroupCache, gID)
}

func (c *Cache) addCategory(dest, bID string, category *Category) {
	if !c.trackAPICalls || category == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.Entries[bID]; !ok {
		return
	}

	c.Entries[bID].CategoryCache[category.ID.String()] = &cacheEntry[Category]{
		Data:        category,
		CreatedAt:   time.Now().UTC(),
		EndpointURL: dest,
	}
}

func (c *Cache) addCategories(dest, bID string, categories []*Category) {
	if !c.trackBulkAPICalls || categories == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.Entries[bID]; !ok {
		return
	}
	for _, cat := range categories {
		if cat == nil {
			continue
		}
		c.Entries[bID].CategoryCache[cat.ID.String()] = &cacheEntry[Category]{
			Data:        cat,
			CreatedAt:   time.Now().UTC(),
			EndpointURL: dest,
		}
	}
}

func (c *Cache) deleteCategory(bID, cID string) {
	if !c.trackAPICalls {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.Entries[bID]; !ok {
		return
	}

	delete(c.Entries[bID].CategoryCache, cID)
}

func (c *Cache) addTxn(dest, bID string, txn *Transaction) {
	if !c.trackAPICalls || txn == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.Entries[bID]; !ok {
		return
	}

	c.Entries[bID].TxnCache[txn.ID.String()] = &cacheEntry[Transaction]{
		Data:        txn,
		CreatedAt:   time.Now().UTC(),
		EndpointURL: dest,
	}
}

func (c *Cache) addTxns(dest, bID string, txns []*Transaction) {
	if !c.trackBulkAPICalls || txns == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.Entries[bID]; !ok {
		return
	}
	for _, t := range txns {
		if t == nil {
			continue
		}
		c.Entries[bID].TxnCache[t.ID.String()] = &cacheEntry[Transaction]{
			Data:        t,
			CreatedAt:   time.Now().UTC(),
			EndpointURL: dest,
		}
	}
}

func (c *Cache) deleteTxn(bID, tID string) {
	if !c.trackAPICalls {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.Entries[bID]; !ok {
		return
	}

	delete(c.Entries[bID].TxnCache, tID)
}

func (c *Cache) addTxnDetails(dest, bID string, txn *TransactionDetail) {
	if !c.trackAPICalls || txn == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.Entries[bID]; !ok {
		return
	}

	c.Entries[bID].TxnDetailsCache[txn.ID.String()] = &cacheEntry[TransactionDetail]{
		Data:        txn,
		CreatedAt:   time.Now().UTC(),
		EndpointURL: dest,
	}
}

func (c *Cache) addTxnsDetails(dest, bID string, txns []*TransactionDetail) {
	if !c.trackBulkAPICalls || txns == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.Entries[bID]; !ok {
		return
	}
	for _, t := range txns {
		if t == nil {
			continue
		}
		c.Entries[bID].TxnDetailsCache[t.ID.String()] = &cacheEntry[TransactionDetail]{
			Data:        t,
			CreatedAt:   time.Now().UTC(),
			EndpointURL: dest,
		}
	}
}

func (c *Cache) deleteTxnsDetails(bID, tID string) {
	if !c.trackAPICalls {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.Entries[bID]; !ok {
		return
	}

	delete(c.Entries[bID].TxnDetailsCache, tID)
}

// Set copies the input entries map of budgetIDs->budgetCaches
// to the cache
func (c *Cache) Set(entries map[string]*budgetCache) {
	c.mu.Lock()
	defer c.mu.Unlock()

	maps.Copy(c.Entries, entries)
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for range ticker.C {
		c.reap()
	}
}

func reapSubcache[T any](m subcache[T], interval time.Duration) {
	for k, entry := range m {
		if time.Since(entry.CreatedAt) > interval {
			delete(m, k)
		}
	}
}

func (c *Cache) reap() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, bCache := range c.Entries {
		if time.Since(bCache.BudgetEntry.CreatedAt) > c.interval {
			delete(c.Entries, k)
		} else {
			reapSubcache(bCache.AccountCache, c.interval)
			reapSubcache(bCache.PayeeCache, c.interval)
			reapSubcache(bCache.GroupCache, c.interval)
			reapSubcache(bCache.CategoryCache, c.interval)
			reapSubcache(bCache.TxnCache, c.interval)
			reapSubcache(bCache.TxnDetailsCache, c.interval)
		}
	}
}

// Clear deletes all cached entries.
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	clear(c.Entries)
}

func NewCache(interval time.Duration) *Cache {
	cache := Cache{
		mu:                &sync.Mutex{},
		interval:          interval,
		trackAPICalls:     true,
		trackBulkAPICalls: true,
		Entries:           map[string]*budgetCache{},
	}

	go cache.reapLoop()

	return &cache
}
