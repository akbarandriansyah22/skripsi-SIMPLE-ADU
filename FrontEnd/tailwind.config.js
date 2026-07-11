export default {
  content: ["./index.html", "./src/**/*.{vue,js,ts,jsx,tsx}"],
  theme: {
    extend: {
      boxShadow: {
        soft: "0 1px 3px rgba(0, 0, 0, 0.08), 0 2px 8px rgba(0, 0, 0, 0.06)",
        card: "0 2px 8px rgba(0, 0, 0, 0.06), 0 4px 16px rgba(0, 0, 0, 0.08)",
        elevated:
          "0 4px 12px rgba(0, 0, 0, 0.08), 0 8px 24px rgba(0, 0, 0, 0.12)",
      },
      borderRadius: {
        soft: "12px",
        card: "16px",
        lg: "20px",
      },
      colors: {
        surface: "#F8FAFC",
        panel: "#F1F5F9",
        navy: "#0F172A",
        primary: "#0B2A4A",
        secondary: "#2563EB",
        success: "#10B981",
        warning: "#F59E0B",
        danger: "#EF4444",
        background: "#F8FAFC",
        card: "#FFFFFF",
        border: "#E2E8F0",
      },
      fontSize: {
        xs: ["12px", { lineHeight: "16px", letterSpacing: "0.3px" }],
        sm: ["13px", { lineHeight: "20px", letterSpacing: "0.2px" }],
        base: ["14px", { lineHeight: "22px", letterSpacing: "0.1px" }],
        lg: ["16px", { lineHeight: "24px", letterSpacing: "0px" }],
        xl: ["18px", { lineHeight: "28px", letterSpacing: "0px" }],
        "2xl": ["20px", { lineHeight: "32px", letterSpacing: "0px" }],
        "3xl": ["28px", { lineHeight: "36px", letterSpacing: "-0.5px" }],
        "4xl": ["36px", { lineHeight: "44px", letterSpacing: "-1px" }],
      },
    },
  },
  plugins: [require("@tailwindcss/forms")],
};
