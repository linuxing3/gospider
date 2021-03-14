import {
  Database,
  PostgresConnector,
  DataTypes,
  Model,
} from "https://deno.land/x/denodb/mod.ts";
import { config } from "https://deno.land/x/dotenv/mod.ts";

class Movie extends Model {
  static table = "movies";
  static timestamps = true;

  static fields = {
    id: { primaryKey: true, autoIncrement: true },
    year: DataTypes.STRING,
    title: DataTypes.STRING,
    subtitle: DataTypes.STRING,
    desc: DataTypes.STRING,
    other: DataTypes.STRING,
    area: DataTypes.STRING,
    tag: DataTypes.STRING,
    star: DataTypes.STRING,
    comment: DataTypes.STRING,
    quote: DataTypes.STRING,
  };
}

const env = config({ safe: true });
export const postOptions = {
  host: env["POST_HOST"],
  username: env["POST_USERNAME"],
  password: env["POST_PASSWORD"],
  port: parseInt(env["POST_PORT"]),
  database: env["POST_DATABASE"],
};
const connection = new PostgresConnector(postOptions);

const postdb = new Database(connection)
postdb.link([Movie]);
await postdb.sync({ drop: true });
await postdb.close();
