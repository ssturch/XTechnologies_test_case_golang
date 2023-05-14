package internal

import (
	"fmt"
	"net/url"
)

type CustomErr struct {
	Tp    string
	Cause string
	Text  string
	Err   error
}

func (ce *CustomErr) Error() string {
	DetectErrType(ce)
	return fmt.Sprintf("Error type: %v | Error text: %v | Error cause: %v\n", ce.Tp, ce.Text, ce.Cause)
}

type RequestError struct {
	Host   string
	Status int
	Text   string
	Err    error
}

func (re *RequestError) Error() string {
	return fmt.Sprintf("Host %v return status %v with text \"%v\"", re.Host, re.Status, re.Text)
}

func (re *RequestError) Unwrap() error {
	return re.Err
}

type ReaderError struct {
	ReadFrom string
	Err      error
}

func (readerr *ReaderError) Error() string {
	return fmt.Sprintf("Read from %v return \"%v\"", readerr.ReadFrom, readerr.Err)
}

func (readerr *ReaderError) Unwrap() error {
	return readerr.Err
}

type DecodeUnmarshallError struct {
	Do    string
	Cause string
	Err   error
}

func (duerr *DecodeUnmarshallError) Error() string {
	return fmt.Sprintf("%v %v caused \"%v\"", duerr.Do, duerr.Cause, duerr.Err)
}

func (duerr *DecodeUnmarshallError) Unwrap() error {
	return duerr.Err
}

type EncodeMarshallError struct {
	Do    string
	Cause string
	Err   error
}

func (emerr *EncodeMarshallError) Error() string {
	return fmt.Sprintf("%v %v caused \"%v\"", emerr.Do, emerr.Cause, emerr.Err)
}

func (emerr *EncodeMarshallError) Unwrap() error {
	return emerr.Err
}

type ParsingError struct {
	From string
	To   string
	Err  error
}

func (perr *ParsingError) Error() string {
	return fmt.Sprintf("Parsing from %v to %v caused \"%v\"", perr.From, perr.To, perr.Err)
}

func (perr *ParsingError) Unwrap() error {
	return perr.Err
}

type DBError struct {
	QueryFunc string
	Err       error
}

func (dberr *DBError) Error() string {
	return fmt.Sprintf("Operation %v caused \"%v\"", dberr.QueryFunc, dberr.Err)
}

func (dberr *DBError) Unwrap() error {
	return dberr.Err
}

func DetectErrType(c *CustomErr) {
	switch c.Err.(type) {
	default:
		fmt.Println("not a model missing error")
	case *url.Error:
		c.Tp = "Internal"
		c.Cause = c.Err.(*url.Error).URL
		c.Text = "Failure to speak HTTP or error caused by server policy"

	case *RequestError:
		c.Tp = "External"
		c.Cause = c.Err.(*RequestError).Host
		c.Text = c.Err.(*RequestError).Error()

	case *ReaderError:
		c.Tp = "Internal Reader"
		c.Cause = c.Err.(*ReaderError).ReadFrom
		c.Text = c.Err.(*ReaderError).Error()

	case *DecodeUnmarshallError:
		c.Tp = "Internal Decode/Unmarshall"
		c.Cause = c.Err.(*DecodeUnmarshallError).Cause
		c.Text = c.Err.(*DecodeUnmarshallError).Error()

	case *EncodeMarshallError:
		c.Tp = "Internal Decode/Unmarshall"
		c.Cause = c.Err.(*DecodeUnmarshallError).Cause
		c.Text = c.Err.(*DecodeUnmarshallError).Error()

	case *ParsingError:
		c.Tp = "Internal Encode/Marshall"
		c.Cause = c.Err.(*ParsingError).From
		c.Text = c.Err.(*ParsingError).Error()

	case *DBError:
		c.Tp = "Internal DB"
		c.Cause = c.Err.(*DBError).QueryFunc
		c.Text = c.Err.(*DBError).Error()
	}

}

// Будущий логгер для ошибок
func ErrorLogger(err error) {
	out := new(CustomErr)
	out.Err = err
	fmt.Println(out.Error())
}
