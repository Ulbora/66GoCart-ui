package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	lg "github.com/Ulbora/Level_Logger"
	conts "github.com/Ulbora/Six910-ui/contentsrv"
	m "github.com/Ulbora/Six910-ui/managers"
	mapi "github.com/Ulbora/Six910-ui/mockapi"
	api "github.com/Ulbora/Six910API-Go"
	ml "github.com/Ulbora/go-mail-sender"
	ds "github.com/Ulbora/json-datastore"
	sdbi "github.com/Ulbora/six910-database-interface"
)

func TestSix910Handler_CustomerLoginPage(t *testing.T) {
	var sh Six910Handler
	var l lg.Logger
	l.LogLevel = lg.AllLevel
	sh.Log = &l
	var sapi mapi.MockAPI
	sapi.SetStoreID(59)

	var man m.Six910Manager
	man.API = &sapi
	sh.API = &sapi
	man.Log = &l
	sh.Manager = man.GetNew()

	sapi.SetRestURL("http://localhost:3002")
	sapi.SetStore("defaultLocalStore", "defaultLocalStore.mydomain.com")
	sapi.SetAPIKey("GDG651GFD66FD16151sss651f651ff65555ddfhjklyy5")
	sh.API = &sapi

	//-----------start mocking------------------

	var prod sdbi.Product
	prod.ID = 2
	prod.Desc = "test"

	var plst []sdbi.Product
	plst = append(plst, prod)
	sapi.MockProductList = &plst

	//-----------end mocking --------

	var c conts.CmsService
	var ds ds.DataStore
	ds.Path = "../contentsrv/testFiles"
	//ds.Delete("books1")
	c.Log = &l
	c.Store = ds.GetNew()

	var ct conts.Content
	ct.Name = "index"
	ct.Author = "ken"
	ct.MetaAuthorName = "ken"
	ct.MetaDesc = "shopping cart index"
	ct.Text = "some book text"
	ct.Title = "the best book ever"
	ct.Visible = true
	res := c.AddContent(&ct)
	fmt.Println("content save: ", res)

	sh.ContentService = c.GetNew()

	var cc ClientCreds
	cc.AuthCodeState = "123"
	sh.ClientCreds = &cc
	sh.ClientCreds.AuthCodeClient = "1"
	sh.OauthHost = "test.com"

	sh.Templates = template.Must(template.ParseFiles("testHtmls/test.html"))

	r, _ := http.NewRequest("POST", "https://test.com", nil)
	w := httptest.NewRecorder()
	s, suc := sh.getUserSession(w, r)
	fmt.Println("suc: ", suc)
	//s.Values["loggedIn"] = true
	s.Save(r, w)
	h := sh.GetNew()
	h.CustomerLoginPage(w, r)
	fmt.Println("code: ", w.Code)

	if w.Code != 200 {
		t.Fail()
	}
}

func TestSix910Handler_CustomerLogin(t *testing.T) {
	var sh Six910Handler
	var l lg.Logger
	l.LogLevel = lg.AllLevel
	sh.Log = &l

	var sapi mapi.MockAPI
	sapi.SetStoreID(59)

	var man m.Six910Manager
	man.API = &sapi
	sh.API = &sapi
	man.Log = &l
	sh.Manager = man.GetNew()

	sapi.SetRestURL("http://localhost:3002")
	sapi.SetStore("defaultLocalStore", "defaultLocalStore.mydomain.com")
	sapi.SetAPIKey("GDG651GFD66FD16151sss651f651ff65555ddfhjklyy5")

	//-----------start mocking------------------
	var user api.UserResponse
	user.Username = "tester123"
	user.Role = customerRole
	user.Enabled = true

	sapi.MockUser = &user

	var ur api.Response
	ur.Success = true
	sapi.MockUpdateUserResp = &ur

	var crt sdbi.Cart
	crt.ID = 4
	sapi.MockCart = &crt

	sh.API = &sapi

	//-----------end mocking --------

	r, _ := http.NewRequest("POST", "/test", strings.NewReader("username=tester123&password=tester"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	//r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h := sh.GetNew()
	h.CustomerLogin(w, r)
	fmt.Println("code: ", w.Code)

	if w.Code != 302 {
		t.Fail()
	}
}

func TestSix910Handler_CustomerLogin2(t *testing.T) {
	var sh Six910Handler
	var l lg.Logger
	l.LogLevel = lg.AllLevel
	sh.Log = &l

	var sapi mapi.MockAPI
	sapi.SetStoreID(59)

	var man m.Six910Manager
	man.API = &sapi
	sh.API = &sapi
	man.Log = &l
	sh.Manager = man.GetNew()

	sapi.SetRestURL("http://localhost:3002")
	sapi.SetStore("defaultLocalStore", "defaultLocalStore.mydomain.com")
	sapi.SetAPIKey("GDG651GFD66FD16151sss651f651ff65555ddfhjklyy5")

	//-----------start mocking------------------
	var user api.UserResponse
	user.Username = "tester123"
	user.Role = customerRole
	user.Enabled = true

	sapi.MockUser = &user

	var ur api.Response
	ur.Success = true
	sapi.MockUpdateUserResp = &ur

	var crt sdbi.Cart
	crt.ID = 4
	sapi.MockCart = &crt

	var cilstp []sdbi.CartItem

	var ctit1 sdbi.CartItem
	ctit1.Quantity = 3
	ctit1.ProductID = 7
	cilstp = append(cilstp, ctit1)

	var ctit2 sdbi.CartItem
	ctit2.Quantity = 4
	ctit2.ProductID = 9
	cilstp = append(cilstp, ctit2)

	var cccs m.CustomerCart
	cccs.Items = &cilstp

	sapi.MockCartItemList = &cilstp

	sh.API = &sapi

	//-----------end mocking --------

	r, _ := http.NewRequest("POST", "/test", strings.NewReader("username=tester123&password=tester"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	//r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s, suc := sh.getUserSession(w, r)
	fmt.Println("suc: ", suc)
	b, _ := json.Marshal(cccs)
	bb := sh.compressObj(b)
	s.Values["customerCart"] = bb
	s.Save(r, w)

	h := sh.GetNew()
	h.CustomerLogin(w, r)
	fmt.Println("code: ", w.Code)

	if w.Code != 302 {
		t.Fail()
	}
}

func TestSix910Handler_CustomerLoginFail(t *testing.T) {
	var sh Six910Handler
	var l lg.Logger
	l.LogLevel = lg.AllLevel
	sh.Log = &l

	var sapi mapi.MockAPI
	sapi.SetStoreID(59)

	var man m.Six910Manager
	man.API = &sapi
	sh.API = &sapi
	man.Log = &l
	sh.Manager = man.GetNew()

	sapi.SetRestURL("http://localhost:3002")
	sapi.SetStore("defaultLocalStore", "defaultLocalStore.mydomain.com")
	sapi.SetAPIKey("GDG651GFD66FD16151sss651f651ff65555ddfhjklyy5")

	//-----------start mocking------------------
	var user api.UserResponse
	user.Username = "tester1234"
	user.Role = customerRole
	user.Enabled = true

	sapi.MockUser = &user

	var ur api.Response
	ur.Success = true
	sapi.MockUpdateUserResp = &ur

	sh.API = &sapi

	//-----------end mocking --------

	r, _ := http.NewRequest("POST", "/test", strings.NewReader("username=tester123&password=tester"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	//r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h := sh.GetNew()
	h.CustomerLogin(w, r)
	fmt.Println("code: ", w.Code)

	if w.Code != 302 {
		t.Fail()
	}
}

func TestSix910Handler_CustomerChangePasswordPage(t *testing.T) {
	var sh Six910Handler
	var l lg.Logger
	l.LogLevel = lg.AllLevel
	sh.Log = &l
	var sapi mapi.MockAPI
	sapi.SetStoreID(59)

	sapi.SetRestURL("http://localhost:3002")
	sapi.SetStore("defaultLocalStore", "defaultLocalStore.mydomain.com")
	sapi.SetAPIKey("GDG651GFD66FD16151sss651f651ff65555ddfhjklyy5")
	sh.API = &sapi

	var man m.Six910Manager
	man.API = &sapi
	sh.API = &sapi
	man.Log = &l
	sh.Manager = man.GetNew()

	//-----------start mocking------------------

	var cusm sdbi.Customer
	//cusm.ID = 5
	cusm.Email = "test@tester.com"
	sapi.MockCustomer = &cusm

	var macres api.Response
	macres.Success = true

	sapi.MockUpdateCustomerResp = &macres

	var mu api.UserResponse
	mu.Enabled = true
	mu.Username = "tester"

	sapi.MockUser = &mu

	var aares api.ResponseID
	aares.Success = true
	aares.ID = 8
	sapi.MockAddAddressRes = &aares

	var aures api.Response
	aures.Success = true

	sapi.MockAddCustomerUserRes = &aures

	//-----------end mocking --------

	var c conts.CmsService
	var ds ds.DataStore
	ds.Path = "../contentsrv/testFiles"
	//ds.Delete("books1")
	c.Log = &l
	c.Store = ds.GetNew()

	var ct conts.Content
	ct.Name = "product"
	ct.Author = "ken"
	ct.MetaAuthorName = "ken"
	ct.MetaDesc = "shopping cart index"
	ct.Text = "some book text"
	ct.Title = "the best book ever"
	ct.Visible = true
	res := c.AddContent(&ct)
	fmt.Println("content save: ", res)

	sh.ContentService = c.GetNew()

	var cc ClientCreds
	cc.AuthCodeState = "123"
	sh.ClientCreds = &cc
	sh.ClientCreds.AuthCodeClient = "1"
	sh.OauthHost = "test.com"

	sh.Templates = template.Must(template.ParseFiles("testHtmls/test.html"))

	r, _ := http.NewRequest("POST", "https://test.com", strings.NewReader("email=bob@bob.com&firstName=tester&"+
		"lastName=testertest&zip=12345&billAddress=123&billCity=dd&billState=rr&billZip=22&"+
		"billCountry=55&shipAddress=444&shipCity=444&shipState=dfg&shipZip=234&shipCountry=55"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	w := httptest.NewRecorder()
	s, suc := sh.getUserSession(w, r)
	fmt.Println("suc: ", suc)
	var cccs m.CustomerCart
	s.Values["userLoggenIn"] = true
	s.Values["customerUser"] = true
	s.Values["customerId"] = int64(55)
	s.Values["customerCart"] = &cccs
	s.Values["username"] = "tester"
	s.Values["password"] = "tester"
	s.Save(r, w)
	h := sh.GetNew()
	h.CustomerChangePasswordPage(w, r)
	fmt.Println("code: ", w.Code)

	if w.Code != 200 {
		t.Fail()
	}
}

func TestSix910Handler_CustomerChangePasswordPageLogin(t *testing.T) {
	var sh Six910Handler
	var l lg.Logger
	l.LogLevel = lg.AllLevel
	sh.Log = &l
	var sapi mapi.MockAPI
	sapi.SetStoreID(59)

	sapi.SetRestURL("http://localhost:3002")
	sapi.SetStore("defaultLocalStore", "defaultLocalStore.mydomain.com")
	sapi.SetAPIKey("GDG651GFD66FD16151sss651f651ff65555ddfhjklyy5")
	sh.API = &sapi

	var man m.Six910Manager
	man.API = &sapi
	sh.API = &sapi
	man.Log = &l
	sh.Manager = man.GetNew()

	//-----------start mocking------------------

	var cusm sdbi.Customer
	//cusm.ID = 5
	cusm.Email = "test@tester.com"
	sapi.MockCustomer = &cusm

	var macres api.Response
	macres.Success = true

	sapi.MockUpdateCustomerResp = &macres

	var mu api.UserResponse
	mu.Enabled = true
	mu.Username = "tester"

	sapi.MockUser = &mu

	var aares api.ResponseID
	aares.Success = true
	aares.ID = 8
	sapi.MockAddAddressRes = &aares

	var aures api.Response
	aures.Success = true

	sapi.MockAddCustomerUserRes = &aures

	//-----------end mocking --------

	var c conts.CmsService
	var ds ds.DataStore
	ds.Path = "../contentsrv/testFiles"
	//ds.Delete("books1")
	c.Log = &l
	c.Store = ds.GetNew()

	var ct conts.Content
	ct.Name = "product"
	ct.Author = "ken"
	ct.MetaAuthorName = "ken"
	ct.MetaDesc = "shopping cart index"
	ct.Text = "some book text"
	ct.Title = "the best book ever"
	ct.Visible = true
	res := c.AddContent(&ct)
	fmt.Println("content save: ", res)

	sh.ContentService = c.GetNew()

	var cc ClientCreds
	cc.AuthCodeState = "123"
	sh.ClientCreds = &cc
	sh.ClientCreds.AuthCodeClient = "1"
	sh.OauthHost = "test.com"

	sh.Templates = template.Must(template.ParseFiles("testHtmls/test.html"))

	r, _ := http.NewRequest("POST", "https://test.com", strings.NewReader("email=bob@bob.com&firstName=tester&"+
		"lastName=testertest&zip=12345&billAddress=123&billCity=dd&billState=rr&billZip=22&"+
		"billCountry=55&shipAddress=444&shipCity=444&shipState=dfg&shipZip=234&shipCountry=55"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	w := httptest.NewRecorder()
	s, suc := sh.getUserSession(w, r)
	fmt.Println("suc: ", suc)
	var cccs m.CustomerCart
	//s.Values["loggedIn"] = true
	s.Values["customerUser"] = true
	s.Values["customerId"] = int64(55)
	s.Values["customerCart"] = &cccs
	s.Values["username"] = "tester"
	s.Values["password"] = "tester"
	s.Save(r, w)
	h := sh.GetNew()
	h.CustomerChangePasswordPage(w, r)
	fmt.Println("code: ", w.Code)

	if w.Code != 302 {
		t.Fail()
	}
}

func TestSix910Handler_CustomerChangePassword(t *testing.T) {
	var sh Six910Handler
	var l lg.Logger
	l.LogLevel = lg.AllLevel
	sh.Log = &l
	var sapi mapi.MockAPI
	sapi.SetStoreID(59)

	sapi.SetRestURL("http://localhost:3002")
	sapi.SetStore("defaultLocalStore", "defaultLocalStore.mydomain.com")
	sapi.SetAPIKey("GDG651GFD66FD16151sss651f651ff65555ddfhjklyy5")
	sh.API = &sapi

	var man m.Six910Manager
	man.API = &sapi
	sh.API = &sapi
	man.Log = &l
	sh.Manager = man.GetNew()

	//-----------start mocking------------------

	var cusm sdbi.Customer
	//cusm.ID = 5
	cusm.Email = "test@tester.com"
	sapi.MockCustomer = &cusm

	var macres api.Response
	macres.Success = true

	sapi.MockUpdateCustomerResp = &macres

	var mu api.UserResponse
	mu.Enabled = true
	mu.Username = "tester"
	mu.Role = customerRole

	sapi.MockUser = &mu

	var aares api.ResponseID
	aares.Success = true
	aares.ID = 8
	sapi.MockAddAddressRes = &aares

	var aures api.Response
	aures.Success = true

	sapi.MockUpdateUserResp = &aures

	//-----------end mocking --------

	var c conts.CmsService
	var ds ds.DataStore
	ds.Path = "../contentsrv/testFiles"
	//ds.Delete("books1")
	c.Log = &l
	c.Store = ds.GetNew()

	var ct conts.Content
	ct.Name = "product"
	ct.Author = "ken"
	ct.MetaAuthorName = "ken"
	ct.MetaDesc = "shopping cart index"
	ct.Text = "some book text"
	ct.Title = "the best book ever"
	ct.Visible = true
	res := c.AddContent(&ct)
	fmt.Println("content save: ", res)

	sh.ContentService = c.GetNew()

	var cc ClientCreds
	cc.AuthCodeState = "123"
	sh.ClientCreds = &cc
	sh.ClientCreds.AuthCodeClient = "1"
	sh.OauthHost = "test.com"

	sh.Templates = template.Must(template.ParseFiles("testHtmls/test.html"))

	r, _ := http.NewRequest("POST", "https://test.com", strings.NewReader("password=tester&oldPassword=tester2&"+
		"lastName=testertest&zip=12345&billAddress=123&billCity=dd&billState=rr&billZip=22&"+
		"billCountry=55&shipAddress=444&shipCity=444&shipState=dfg&shipZip=234&shipCountry=55"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	w := httptest.NewRecorder()
	s, suc := sh.getUserSession(w, r)
	fmt.Println("suc: ", suc)
	var cccs m.CustomerCart
	s.Values["userLoggenIn"] = true
	s.Values["customerUser"] = true
	s.Values["customerId"] = int64(55)
	s.Values["customerCart"] = &cccs
	s.Values["username"] = "tester"
	s.Values["password"] = "tester"
	s.Save(r, w)
	h := sh.GetNew()
	h.CustomerChangePassword(w, r)
	fmt.Println("code: ", w.Code)

	if w.Code != 302 {
		t.Fail()
	}
}

func TestSix910Handler_CustomerChangePasswordLogin(t *testing.T) {
	var sh Six910Handler
	var l lg.Logger
	l.LogLevel = lg.AllLevel
	sh.Log = &l
	var sapi mapi.MockAPI
	sapi.SetStoreID(59)

	sapi.SetRestURL("http://localhost:3002")
	sapi.SetStore("defaultLocalStore", "defaultLocalStore.mydomain.com")
	sapi.SetAPIKey("GDG651GFD66FD16151sss651f651ff65555ddfhjklyy5")
	sh.API = &sapi

	var man m.Six910Manager
	man.API = &sapi
	sh.API = &sapi
	man.Log = &l
	sh.Manager = man.GetNew()

	//-----------start mocking------------------

	var cusm sdbi.Customer
	//cusm.ID = 5
	cusm.Email = "test@tester.com"
	sapi.MockCustomer = &cusm

	var macres api.Response
	macres.Success = true

	sapi.MockUpdateCustomerResp = &macres

	var mu api.UserResponse
	mu.Enabled = true
	mu.Username = "tester"

	sapi.MockUser = &mu

	var aares api.ResponseID
	aares.Success = true
	aares.ID = 8
	sapi.MockAddAddressRes = &aares

	var aures api.Response
	aures.Success = true

	sapi.MockUpdateUserResp = &aures

	//-----------end mocking --------

	var c conts.CmsService
	var ds ds.DataStore
	ds.Path = "../contentsrv/testFiles"
	//ds.Delete("books1")
	c.Log = &l
	c.Store = ds.GetNew()

	var ct conts.Content
	ct.Name = "product"
	ct.Author = "ken"
	ct.MetaAuthorName = "ken"
	ct.MetaDesc = "shopping cart index"
	ct.Text = "some book text"
	ct.Title = "the best book ever"
	ct.Visible = true
	res := c.AddContent(&ct)
	fmt.Println("content save: ", res)

	sh.ContentService = c.GetNew()

	var cc ClientCreds
	cc.AuthCodeState = "123"
	sh.ClientCreds = &cc
	sh.ClientCreds.AuthCodeClient = "1"
	sh.OauthHost = "test.com"

	sh.Templates = template.Must(template.ParseFiles("testHtmls/test.html"))

	r, _ := http.NewRequest("POST", "https://test.com", strings.NewReader("email=bob@bob.com&firstName=tester&"+
		"lastName=testertest&zip=12345&billAddress=123&billCity=dd&billState=rr&billZip=22&"+
		"billCountry=55&shipAddress=444&shipCity=444&shipState=dfg&shipZip=234&shipCountry=55"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	w := httptest.NewRecorder()
	s, suc := sh.getUserSession(w, r)
	fmt.Println("suc: ", suc)
	var cccs m.CustomerCart
	//s.Values["loggedIn"] = true
	s.Values["customerUser"] = true
	s.Values["customerId"] = int64(55)
	s.Values["customerCart"] = &cccs
	s.Values["username"] = "tester"
	s.Values["password"] = "tester"
	s.Save(r, w)
	h := sh.GetNew()
	h.CustomerChangePassword(w, r)
	fmt.Println("code: ", w.Code)

	if w.Code != 302 {
		t.Fail()
	}
}

func TestSix910Handler_CustomerChangePasswordFail(t *testing.T) {
	var sh Six910Handler
	var l lg.Logger
	l.LogLevel = lg.AllLevel
	sh.Log = &l
	var sapi mapi.MockAPI
	sapi.SetStoreID(59)

	sapi.SetRestURL("http://localhost:3002")
	sapi.SetStore("defaultLocalStore", "defaultLocalStore.mydomain.com")
	sapi.SetAPIKey("GDG651GFD66FD16151sss651f651ff65555ddfhjklyy5")
	sh.API = &sapi

	var man m.Six910Manager
	man.API = &sapi
	sh.API = &sapi
	man.Log = &l
	sh.Manager = man.GetNew()

	//-----------start mocking------------------

	var cusm sdbi.Customer
	//cusm.ID = 5
	cusm.Email = "test@tester.com"
	sapi.MockCustomer = &cusm

	var macres api.Response
	macres.Success = true

	sapi.MockUpdateCustomerResp = &macres

	var mu api.UserResponse
	mu.Enabled = true
	mu.Username = "tester"

	sapi.MockUser = &mu

	var aares api.ResponseID
	aares.Success = true
	aares.ID = 8
	sapi.MockAddAddressRes = &aares

	var aures api.Response
	//aures.Success = true

	sapi.MockUpdateUserResp = &aures

	//-----------end mocking --------

	var c conts.CmsService
	var ds ds.DataStore
	ds.Path = "../contentsrv/testFiles"
	//ds.Delete("books1")
	c.Log = &l
	c.Store = ds.GetNew()

	var ct conts.Content
	ct.Name = "product"
	ct.Author = "ken"
	ct.MetaAuthorName = "ken"
	ct.MetaDesc = "shopping cart index"
	ct.Text = "some book text"
	ct.Title = "the best book ever"
	ct.Visible = true
	res := c.AddContent(&ct)
	fmt.Println("content save: ", res)

	sh.ContentService = c.GetNew()

	var cc ClientCreds
	cc.AuthCodeState = "123"
	sh.ClientCreds = &cc
	sh.ClientCreds.AuthCodeClient = "1"
	sh.OauthHost = "test.com"

	sh.Templates = template.Must(template.ParseFiles("testHtmls/test.html"))

	r, _ := http.NewRequest("POST", "https://test.com", strings.NewReader("email=bob@bob.com&firstName=tester&"+
		"lastName=testertest&zip=12345&billAddress=123&billCity=dd&billState=rr&billZip=22&"+
		"billCountry=55&shipAddress=444&shipCity=444&shipState=dfg&shipZip=234&shipCountry=55"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	w := httptest.NewRecorder()
	s, suc := sh.getUserSession(w, r)
	fmt.Println("suc: ", suc)
	var cccs m.CustomerCart
	s.Values["userLoggenIn"] = true
	s.Values["customerUser"] = true
	s.Values["customerId"] = int64(55)
	s.Values["customerCart"] = &cccs
	s.Values["username"] = "tester"
	s.Values["password"] = "tester"
	s.Save(r, w)
	h := sh.GetNew()
	h.CustomerChangePassword(w, r)
	fmt.Println("code: ", w.Code)

	if w.Code != 302 {
		t.Fail()
	}
}

func TestSix910Handler_CustomerLogout(t *testing.T) {
	var sh Six910Handler
	//h.TokenMap = make(map[string]*oauth2.Token)

	// var mTkn oauth2.Token
	// mTkn.AccessToken = "45ffffff"

	// var mockAcTkn oauth2.MockAuthCodeToken
	// mockAcTkn.MockToken = &mTkn

	// h.AuthToken = &mockAcTkn

	var l lg.Logger
	l.LogLevel = lg.AllLevel
	sh.Log = &l
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "https://test.com", nil)
	s, suc := sh.getUserSession(w, r)
	fmt.Println("suc: ", suc)

	s.Values["accessTokenKey"] = "123"
	s.Values["userLoggenIn"] = true
	s.Values["customerUser"] = true

	s.Save(r, w)
	h := sh.GetNew()
	h.CustomerLogout(w, r)
	fmt.Println("code: ", w.Code)
	if w.Code != 302 {
		t.Fail()
	}
}

func TestSix910Handler_CustomerResetPasswordPage(t *testing.T) {
	var sh Six910Handler
	var l lg.Logger
	l.LogLevel = lg.AllLevel
	sh.Log = &l
	var sapi mapi.MockAPI
	sapi.SetStoreID(59)

	var man m.Six910Manager
	man.API = &sapi
	sh.API = &sapi
	man.Log = &l
	sh.Manager = man.GetNew()

	sapi.SetRestURL("http://localhost:3002")
	sapi.SetStore("defaultLocalStore", "defaultLocalStore.mydomain.com")
	sapi.SetAPIKey("GDG651GFD66FD16151sss651f651ff65555ddfhjklyy5")
	sh.API = &sapi

	//-----------start mocking------------------

	var prod sdbi.Product
	prod.ID = 2
	prod.Desc = "test"

	var plst []sdbi.Product
	plst = append(plst, prod)
	sapi.MockProductList = &plst

	//-----------end mocking --------

	var c conts.CmsService
	var ds ds.DataStore
	ds.Path = "../contentsrv/testFiles"
	//ds.Delete("books1")
	c.Log = &l
	c.Store = ds.GetNew()

	var ct conts.Content
	ct.Name = "index"
	ct.Author = "ken"
	ct.MetaAuthorName = "ken"
	ct.MetaDesc = "shopping cart index"
	ct.Text = "some book text"
	ct.Title = "the best book ever"
	ct.Visible = true
	res := c.AddContent(&ct)
	fmt.Println("content save: ", res)

	sh.ContentService = c.GetNew()

	var cc ClientCreds
	cc.AuthCodeState = "123"
	sh.ClientCreds = &cc
	sh.ClientCreds.AuthCodeClient = "1"
	sh.OauthHost = "test.com"

	sh.Templates = template.Must(template.ParseFiles("testHtmls/test.html"))

	r, _ := http.NewRequest("POST", "https://test.com", nil)
	w := httptest.NewRecorder()
	s, suc := sh.getUserSession(w, r)
	fmt.Println("suc: ", suc)
	//s.Values["loggedIn"] = true
	s.Save(r, w)
	h := sh.GetNew()
	h.CustomerResetPasswordPage(w, r)
	fmt.Println("code: ", w.Code)

	if w.Code != 200 {
		t.Fail()
	}
}

func TestSix910Handler_CustomerResetPassword(t *testing.T) {
	var sh Six910Handler
	var l lg.Logger
	l.LogLevel = lg.AllLevel
	sh.Log = &l

	var mss ml.MockSecureSender
	mss.MockSuccess = true
	sh.MailSender = mss.GetNew()
	sh.MailSenderAddress = "test@test.com"

	var sapi mapi.MockAPI
	sapi.SetStoreID(59)

	var man m.Six910Manager
	man.API = &sapi
	sh.API = &sapi
	man.Log = &l
	sh.Manager = man.GetNew()

	sapi.SetRestURL("http://localhost:3002")
	sapi.SetStore("defaultLocalStore", "defaultLocalStore.mydomain.com")
	sapi.SetAPIKey("GDG651GFD66FD16151sss651f651ff65555ddfhjklyy5")

	//-----------start mocking------------------
	var user api.UserResponse
	user.Username = "tester123"
	user.Role = customerRole
	user.Enabled = true

	sapi.MockUser = &user

	var ur api.CustomerPasswordResponse
	ur.Success = true
	ur.Username = "tester123"
	ur.Password = "123"
	sapi.MockCustomerPasswordResp = &ur

	var crt sdbi.Cart
	crt.ID = 4
	sapi.MockCart = &crt

	sh.API = &sapi

	//-----------end mocking --------

	r, _ := http.NewRequest("POST", "/test", strings.NewReader("username=tester123"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	//r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h := sh.GetNew()
	h.CustomerResetPassword(w, r)
	fmt.Println("code: ", w.Code)

	if w.Code != 302 {
		t.Fail()
	}
}
