import { parse } from "jsr:@std/yaml";
import { Project } from "../models/internal-models/models.ts";
import { Project as PublicProject} from "../models/public-models/public-models.ts";

export default async function parseYamlFile(filePath: string) {
    const yml = await Deno.readTextFile(filePath);
    const data = parse(yml) as PublicProject;
    return data;
}
