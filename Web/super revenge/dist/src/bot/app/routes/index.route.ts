import express from "express";
import { renderIndex,
    //  getToken, 
     startBot } from "../controllers/index.controller";

const router = express.Router();

router.get("/", renderIndex);
// router.get("/csrf-token", getToken);
router.post("/", startBot);

export default router;
