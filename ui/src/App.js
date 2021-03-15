import {MDBContainer, MDBNavbar, MDBNavbarBrand,} from "mdbreact";
import {BrowserRouter, Route, Routes} from "react-router-dom";
import Dashboard from "./Dashboard";
import Login from "./Login";


function App() {
    return (
        <div className="App">
            <MDBNavbar color="indigo" dark expand="md">
                <MDBNavbarBrand>
                    <strong className="white-text">WG VPN Server</strong>
                </MDBNavbarBrand>
            </MDBNavbar>
            <MDBContainer>
                <BrowserRouter>
                    <Routes>
                        <Route path={"/dashboard"} element={<Dashboard />} />
                        <Route path={"/"} element={<Login />} />
                    </Routes>
                </BrowserRouter>
            </MDBContainer>
        </div>
    );
}

export default App;


