package server


func (srv *HTTPServer) addRoutes(){
  srv.mux.GET("/testport", handlers.TestPort)
}
