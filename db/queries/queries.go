package queries

const InsertUserQuery = `INSERT INTO public.user (username, password) VALUES ($1, $2)`

const AuthenticateUserQuery = `SELECT username FROM public.user WHERE username=$1 AND password=$2`
