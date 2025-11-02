import { generateService } from "@umijs/openapi";

generateService({
  requestLibPath: "import request from '@/utils/request'",
  schemaPath: "http://127.0.0.1:4523/export/openapi/6?version=3.0",
  serversPath: "./src/api",
  projectName: "gen",
});
