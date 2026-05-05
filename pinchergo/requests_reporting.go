package pinchergo

import (
	"net/http"
)

func (c *Client) BudgetReport(bID, mID string) (report *MonthReport, err error) {
	endpoint := EndpointBudgetMonth(bID, mID)
	err = c.Request(http.MethodGet, endpoint, nil, &report)
	return report, err
}

type categoryReportsContainer struct {
	CategoryReports []*CategoryReport `json:"data"`
}

func (c *Client) BudgetCategoryReports(bID, mID string) (categories []*CategoryReport, err error) {
	endpoint := EndpointBudgetMonthCategories(bID, mID)
	var container categoryReportsContainer
	err = c.Request(http.MethodGet, endpoint, nil, &container)
	return container.CategoryReports, err
}

type groupReportsContainer struct {
	GroupReports []*GroupReport `json:"data"`
}

func (c *Client) BudgetGroupReports(bID, mID string) (groups []*GroupReport, err error) {
	endpoint := EndpointBudgetMonthCategories(bID, mID)
	var container groupReportsContainer
	err = c.Request(http.MethodGet, endpoint, nil, &container)
	return container.GroupReports, err
}
