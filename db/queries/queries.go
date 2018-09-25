package queries

const InsertUserQuery = `INSERT INTO public.user (username, password) VALUES ($1, $2)`

const AuthenticateUserQuery = `SELECT username FROM public.user WHERE username=$1 AND password=$2`

const InsertTokenQuery = `INSERT INTO public.token (username, token) VALUES ($1, $2)`

const GetTokensQuery = `SELECT token, inserted_date FROM public.token WHERE username=$1`

const DeleteTokenQuery = `DELETE FROM token WHERE token=$1`
