package httpclient

import(
	"errors"
	"net/url"
	"strings"
)
// convert host:port:username:pw to -> http://user:pw@host:port
// if already like a URL has :// then returns it as no change

func ProxyURLConvert(s string)(string , error){
	s = strings.TrimSpace(s)
	if s == ""{
		return "", nil
	}

	// already a URL
	if strings.Contains(s, "://"){
		_, err := url.Parse(s)
		if err != nil{
			return "", err
		}

		return s , nil
	}

	parts := strings.Split(s, ":")
	if len(parts) != 4{
		return "", errors.New("invalid format")
	}

	host := parts[0]
	port := parts[1]
	user := parts[2]
	pw   := parts[3]

	u := &url.URL{
		Scheme: "http",
		Host: host + ":" + port,
		User: url.UserPassword(user, pw), 
	}
	return u.String(), nil
}