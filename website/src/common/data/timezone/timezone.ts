import timezones from "./timezones.json";

export type TimeZone {
  label: string;
  tzCode: string;
  name: string;
  utc: string;
}

export default timezones as TimeZone[];
