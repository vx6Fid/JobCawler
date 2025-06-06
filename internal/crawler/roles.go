package crawler

import "strings"

func IsRoleAllowed(role string) bool {
	return AllowedRoles[strings.ToLower(role)]
}

var AllowedRoles = map[string]bool{
	"":                                true, // Allow empty role for general queries
	"backend":                         true,
	"developer":                       true,
	"engineer":                        true,
	"programmer":                      true,
	"devops":                          true,
	"frontend":                        true,
	"data scientist":                  true,
	"mobile developer":                true,
	"fullstack engineer":              true,
	"product manager":                 true,
	"ui/ux designer":                  true,
	"software engineer":               true,
	"system architect":                true,
	"database administrator":          true,
	"cloud engineer":                  true,
	"machine learning engineer":       true,
	"security engineer":               true,
	"network engineer":                true,
	"site reliability engineer":       true,
	"game developer":                  true,
	"business analyst":                true,
	"quality assurance engineer":      true,
	"technical writer":                true,
	"data analyst":                    true,
	"web developer":                   true,
	"research scientist":              true,
	"blockchain developer":            true,
	"AI engineer":                     true,
	"IoT developer":                   true,
	"test automation engineer":        true,
	"scrum master":                    true,
	"product owner":                   true,
	"sales engineer":                  true,
	"application support analyst":     true,
	"ERP consultant":                  true,
	"digital marketing specialist":    true,
	"business intelligence developer": true,
	"data engineer":                   true,
	"technical support engineer":      true,
	"mobile application tester":       true,
	"user researcher":                 true,
	"content strategist":              true,
	"web designer":                    true,
	"IT project manager":              true,
	"system administrator":            true,
	"help desk technician":            true,
	"e-commerce specialist":           true,
	"video game tester":               true,
	"data visualization expert":       true,
	"RPA developer":                   true,
	"chatbot developer":               true,
	"digital product designer":        true,
	"IT consultant":                   true,
	"software development manager":    true,
}
