import { Project } from "../models/internal-models/models.ts";
import parseYamlFile from "./yaml.ts";

export default class Client {
    async deploy(options: { file: string }) {
        const { file } = options;
        console.log("File: ", file);

        const project: Project = await parseYamlFile(file);
        console.log(project);
        console.log(project.project)
    }
}