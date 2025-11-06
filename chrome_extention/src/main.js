import {createApp} from "vue";
import App from "./APP.vue";
import "./index.css"

import "@varlet/ui/es/style";
import Varlet from '@varlet/ui'

createApp(App).use(Varlet).mount("#app");

