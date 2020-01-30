import Vue from "vue";
import Vuetify from "vuetify/lib";
import Logo from "../icons/Logo";
import Skull from "../icons/Skull";
import Stopwatch from "../icons/Stopwatch";
import Gem from "../icons/Gem";
import Dagger from "../icons/Dagger";
import Target from "../icons/Target";
import Splat from "../icons/Splat";

Vue.use(Vuetify);

const DDSTATS_ICONS = {
  logo: {
    component: Logo
  },
  skull: {
    component: Skull
  },
  stopwatch: {
    component: Stopwatch
  },
  gem: {
    component: Gem
  },
  dagger: {
    component: Dagger
  },
  target: {
    component: Target
  },
  splat: {
    component: Splat
  }
};

export default new Vuetify({
  icons: {
    values: DDSTATS_ICONS
  },
  theme: {
    options: {
      customProperties: true
    },
    themes: {
      light: {
        logo: "#34302e",
        background: "#f6f2ee",
        // background: "#fffefc",
        header: "#fffefc",
        // header: "#f6f2ee",
        footer: "#EBE7E4",
        primary: "#34302e",
        secondary: "#424242",
        accent: "#c33409",
        error: "#FF5252",
        info: "#c33409",
        success: "#4CAF50",
        warning: "#FFC107",
        highlight: "#fffefc",
        deselected: "#716f6d"
      },
      dark: {
        logo: "#c33409",
        background: "#242121",
        highlight: "#34302e",
        primary: "#f6f2ee",
        footer: "#1F1F1F",
        accent: "#c33409",
        deselected: "#a7a6a6"
      }
    }
  }
});
