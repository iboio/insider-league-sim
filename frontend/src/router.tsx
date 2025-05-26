import {Routes, Route} from "react-router-dom";
import Home from "@/pages/home.tsx";
import League from "@/pages/league.tsx";

function Router() {
    return (
        <Routes>
            <Route path={"/"} element={<Home/>}></Route>
            <Route path={"/league/:leagueId"} element={<League/>}></Route>
        </Routes>
    )
}

export default Router