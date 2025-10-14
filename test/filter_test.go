package test

import (
	"regexp"
	"testing"
)

type filter struct {
	filter   string
	goodurls []string
	badurls  []string
}

func TestFilters(t *testing.T) {
	tests := []filter{
		{
			filter: `https://djinni\.co/jobs/\d+-\S+`,
			goodurls: []string{
				"https://djinni.co/jobs/760752-qa-automation-trainee-irc272847/",
				"https://djinni.co/jobs/775566-lead-qa-engineer-general-qa/",
				"https://djinni.co/jobs/775623-middle-go-developer-rozetkapay",
			},
			badurls: []string{
				"https://djinni.co/jobs/company-evo/",
				"https://djinni.co/jobs/company-promova/",
			},
		},
		{
			filter: `https://www\.globallogic\.com/ua/careers/\S+-irc\d+`,
			goodurls: []string{
				"https://www.globallogic.com/ua/careers/engineering-director-irc275431/",
				"https://www.globallogic.com/ua/careers/python-technical-lead-irc272738/",
				"https://www.globallogic.com/ua/careers/strong-middle-golang-developer-irc273264/",
			},
			badurls: []string{
				"https://www.globallogic.com/ua/careers/key-projects/",
				"https://www.globallogic.com/ua/careers/why-globallogic/evp_page/",
			},
		},
	}

	for _, test := range tests {
		t.Run("testname", func(t *testing.T) {
			r, err := regexp.Compile(test.filter)
			if err != nil {
				t.Fatal("regexp.Compile")
			} else {
				for _, u := range test.goodurls {
					if !r.MatchString(u) {
						t.Fatal("regexp.MatchString goodurls", u)
					}
				}
				for _, u := range test.badurls {
					if r.MatchString(u) {
						t.Fatal("regexp.MatchString badurls", u)
					}
				}
			}
		})

	}
}
