package errorsStorage

const Settings = 1
const Message = 2
const Auth = 3
const Source = 4
const Destination = 5

const GeneralError = 1
const AppSettingsError = 2
const OptionsError = 3
const InstanceSettingsError = 4
const ApiError = 5
const CredentialsError = 6
const AuthError = 7
const FileError = 8
const FolderError = 9
const NotFoundError = 10
const DownloadError = 11
const UploadError = 12
const EncryptionError = 13
const AccessSharingError = 14
const DeletionError = 15

type Error struct {
	text  string
	scope int
	code  int

	ErrorContract
}

type ErrorContract interface {
	SetError(text string)
	SetCode(code int)
	Code() int
	SetScope(scope int)
	Scope() int

	error
}

func (e *Error) SetError(text string) {
	e.text = text
}

func (e *Error) Error() string {
	return e.text
}

func (e *Error) SetCode(code int) {
	e.code = code
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) SetScope(scope int) {
	e.scope = scope
}

func (e *Error) Scope() int {
	return e.scope
}

func NewError(text string, scope int, code int) *Error {
	return &Error{
		text:  text,
		scope: scope,
		code:  code,
	}
}
