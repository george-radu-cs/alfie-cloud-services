import dotenv from "dotenv";
const env = dotenv.config({
  path: ".env",
  override: true,
});

export const GetIntEnvVar = (name: string): number => {
  const value = parseInt(process.env[name] as string, 10);
  if (isNaN(value)) {
    throw new Error(
      `Error while trying to parse env variable ${name} as Int. Raw value: ${process.env[name]}`
    );
  }
  return value;
};

export const GetStringEnvVar = (name: string): string => {
  const value = process.env[name];
  if (value === undefined) {
    throw new Error(`Env variable ${name} is not defined`);
  }
  return value;
};

export const GetBooleanEnvVar = (name: string): boolean => {
  const value = process.env[name];
  if (value === undefined) {
    throw new Error(`Env variable ${name} is not defined`);
  }
  if (value === "true") {
    return true;
  } else if (value === "false") {
    return false;
  } else {
    throw new Error(
      `Error while trying to parse env variable ${name} as Boolean. Raw value: ${process.env[name]}`
    );
  }
};

export default env;
