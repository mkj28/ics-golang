package ics_test

import (
	"fmt"
	"github.com/PuloV/ics-golang"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestNewParser(t *testing.T) {
	parser := ics.New()
	rType := fmt.Sprintf("%v", reflect.TypeOf(parser))
	if rType != "*ics.Parser" {
		t.Errorf("Failed to create a Parser !")
	}
}

func TestNewParserChans(t *testing.T) {
	parser := ics.New()
	input := parser.GetInputChan()
	output := parser.GetOutputChan()

	rType := fmt.Sprintf("%v", reflect.TypeOf(input))

	if rType != "chan string" {
		t.Errorf("Failed to create a input chan ! Received : Type %s Value %s", rType, input)
	}

	rType = fmt.Sprintf("%v", reflect.TypeOf(output))
	if rType != "chan *ics.Event" {
		t.Errorf("Failed to create a output chan! Received : Type %s Value %s", rType, output)
	}
}

func TestParsing0Calendars(t *testing.T) {
	parser := ics.New()
	parser.Wait()

	parseErrors, err := parser.GetErrors()

	if err != nil {
		t.Errorf("Failed to wait the parse of the calendars ( %s ) \n", err)
	}
	for i, pErr := range parseErrors {
		t.Errorf("Parsing Error №%d : %s  \n", i, pErr)
	}
}

func TestParsing1Calendars(t *testing.T) {
	parser := ics.New()
	input := parser.GetInputChan()
	input <- "testCalendars/2eventsCal.ics"
	parser.Wait()

	parseErrors, err := parser.GetErrors()

	if err != nil {
		t.Errorf("Failed to wait the parse of the calendars ( %s ) \n", err)
	}
	for i, pErr := range parseErrors {
		t.Errorf("Parsing Error №%d : %s  \n", i, pErr)
	}

	calendars, errCal := parser.GetCalendars()

	if errCal != nil {
		t.Errorf("Failed to get calendars ( %s ) \n", errCal)
	}

	if len(calendars) != 1 {
		t.Errorf("Expected 1 calendar , found %d calendars \n", len(calendars))
	}

}

func TestParsing2Calendars(t *testing.T) {
	parser := ics.New()
	input := parser.GetInputChan()
	input <- "testCalendars/2eventsCal.ics"
	input <- "testCalendars/3eventsNoAttendee.ics"
	parser.Wait()

	parseErrors, err := parser.GetErrors()

	if err != nil {
		t.Errorf("Failed to wait the parse of the calendars ( %s ) \n", err)
	}
	for i, pErr := range parseErrors {
		t.Errorf("Parsing Error №%d : %s  \n", i, pErr)
	}

	calendars, errCal := parser.GetCalendars()

	if errCal != nil {
		t.Errorf("Failed to get calendars ( %s ) \n", errCal)
	}

	if len(calendars) != 2 {
		t.Errorf("Expected 1 calendar , found %d calendars \n", len(calendars))
	}

}

func TestParsingNotExistingCalendar(t *testing.T) {
	parser := ics.New()
	input := parser.GetInputChan()
	input <- "testCalendars/notFound.ics"
	parser.Wait()

	parseErrors, err := parser.GetErrors()

	if err != nil {
		t.Errorf("Failed to wait the parse of the calendars ( %s ) \n", err)
	}
	if len(parseErrors) != 1 {
		t.Errorf("Expected 1 error , found %d in :\n  %#v  \n", len(parseErrors), parseErrors)
	}

}

func TestParsingNotExistingAndExistingCalendars(t *testing.T) {
	parser := ics.New()
	input := parser.GetInputChan()
	input <- "testCalendars/3eventsNoAttendee.ics"
	input <- "testCalendars/notFound.ics"
	parser.Wait()

	parseErrors, err := parser.GetErrors()

	if err != nil {
		t.Errorf("Failed to wait the parse of the calendars ( %s ) \n", err)
	}
	if len(parseErrors) != 1 {
		t.Errorf("Expected 1 error , found %d in :\n  %#v  \n", len(parseErrors), parseErrors)
	}

	calendars, errCal := parser.GetCalendars()

	if errCal != nil {
		t.Errorf("Failed to get calendars ( %s ) \n", errCal)
	}

	if len(calendars) != 1 {
		t.Errorf("Expected 1 calendar , found %d calendars \n", len(calendars))
	}

}
func TestParsingWrongCalendarUrls(t *testing.T) {
	parser := ics.New()
	input := parser.GetInputChan()
	input <- "http://localhost/goTestFails"
	parser.Wait()

	parseErrors, err := parser.GetErrors()

	if err != nil {
		t.Errorf("Failed to wait the parse of the calendars ( %s ) \n", err)
	}
	if len(parseErrors) != 1 {
		t.Errorf("Expected 1 error , found %d in :\n  %#v  \n", len(parseErrors), parseErrors)
	}

	calendars, errCal := parser.GetCalendars()

	if errCal != nil {
		t.Errorf("Failed to get calendars ( %s ) \n", errCal)
	}

	if len(calendars) != 0 {
		t.Errorf("Expected 0 calendar , found %d calendars \n", len(calendars))
	}

}

func TestCreatingTempDir(t *testing.T) {

	ics.FilePath = "testingTempDir/"
	parser := ics.New()
	input := parser.GetInputChan()
	input <- "https://www.google.com/calendar/ical/yordanpulov%40gmail.com/private-81525ac0eb14cdc2e858c15e1b296a1c/basic.ics"
	parser.Wait()
	_, err := os.Stat(ics.FilePath)
	if err != nil {
		t.Errorf("Failed to create %s  \n", ics.FilePath)
	}
	// remove the new dir
	os.Remove(ics.FilePath)
	// return the var to default
	ics.FilePath = "tmp/"
}

func TestCalendarInfo(t *testing.T) {
	parser := ics.New()
	input := parser.GetInputChan()
	input <- "testCalendars/2eventsCal.ics"
	parser.Wait()

	parseErrors, err := parser.GetErrors()

	if err != nil {
		t.Errorf("Failed to wait the parse of the calendars ( %s ) \n", err)
	}
	if len(parseErrors) != 0 {
		t.Errorf("Expected 0 error , found %d in :\n  %#v  \n", len(parseErrors), parseErrors)
	}

	calendars, errCal := parser.GetCalendars()

	if errCal != nil {
		t.Errorf("Failed to get calendars ( %s ) \n", errCal)
	}

	if len(calendars) != 1 {
		t.Errorf("Expected 1 calendar , found %d calendars \n", len(calendars))
		return
	}

	calendar := calendars[0]

	if calendar.GetName() != "2 Events Cal" {
		t.Errorf("Expected name '%s' calendar , got '%s' calendars \n", "2 Events Cal", calendar.GetName())
	}

	if calendar.GetDesc() != "The cal has 2 events(1st with attendees and second without)" {
		t.Errorf("Expected description '%s' calendar , got '%s' calendars \n", "The cal has 2 events(1st with attendees and second without)", calendar.GetDesc())
	}

	if calendar.GetVersion() != 2.0 {
		t.Errorf("Expected version %s calendar , got %s calendars \n", 2.0, calendar.GetVersion())
	}

	events := calendar.GetEvents()
	if len(events) != 2 {
		t.Errorf("Expected  %s events in calendar , got %s events \n", 2, len(events))
	}

	eventsByDates := calendar.GetEventsByDates()
	if len(eventsByDates) != 2 {
		t.Errorf("Expected  %s events in calendar , got %s events \n", 2, len(eventsByDates))
	}

	geometryExamIcsFormat, errICS := time.Parse(ics.IcsFormat, "20140616T060000Z")
	if err != nil {
		t.Errorf("(ics time format) Unexpected error %s \n", errICS)
	}

	geometryExamYmdHis, errYMD := time.Parse(ics.YmdHis, "2014-06-16 06:00:00")
	if err != nil {
		t.Errorf("(YmdHis time format) Unexpected error %s \n", errYMD)
	}
	eventsByDate, err := calendar.GetEventsByDate(geometryExamIcsFormat)
	if err != nil {
		t.Errorf("(ics time format) Unexpected error %s \n", err)
	}
	if len(eventsByDate) != 1 {
		t.Errorf("(ics time format) Expected  %s events in calendar for the date 2014-06-16 , got %s events \n", 1, len(eventsByDate))
	}

	eventsByDate, err = calendar.GetEventsByDate(geometryExamYmdHis)
	if err != nil {
		t.Errorf("(YmdHis time format) Unexpected error %s \n", err)
	}
	if len(eventsByDate) != 1 {
		t.Errorf("(YmdHis time format) Expected  %s events in calendar for the date 2014-06-16 , got %s events \n", 1, len(eventsByDate))
	}

}

func TestCalendarEvents(t *testing.T) {
	parser := ics.New()
	input := parser.GetInputChan()
	input <- "testCalendars/2eventsCal.ics"
	parser.Wait()

	parseErrors, err := parser.GetErrors()

	if err != nil {
		t.Errorf("Failed to wait the parse of the calendars ( %s ) \n", err)
	}
	if len(parseErrors) != 0 {
		t.Errorf("Expected 0 error , found %d in :\n  %#v  \n", len(parseErrors), parseErrors)
	}

	calendars, errCal := parser.GetCalendars()

	if errCal != nil {
		t.Errorf("Failed to get calendars ( %s ) \n", errCal)
	}

	if len(calendars) != 1 {
		t.Errorf("Expected 1 calendar , found %d calendars \n", len(calendars))
		return
	}

	calendar := calendars[0]
	event, err := calendar.GetEventByImportedID("btb9tnpcnd4ng9rn31rdo0irn8@google.com")
	if err != nil {
		t.Errorf("Failed to get event by id with error %s \n", err)
	}

	start, _ := time.Parse(ics.IcsFormat, "20140714T100000Z")
	end, _ := time.Parse(ics.IcsFormat, "20140714T110000Z")

	if event.GetStart() != start {
		t.Errorf("Expected start %s , found %s  \n", start, event.GetStart())
	}
	if event.GetEnd() != end {
		t.Errorf("Expected start %s , found %s  \n", end, event.GetEnd())
	}
}
