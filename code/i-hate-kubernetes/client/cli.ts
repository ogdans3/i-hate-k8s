import { Command } from "@commander-js/extra-typings";
import Client from "./client.ts";
const program = new Command();

program
    .name("I hate kubernetes")
    .description("Boy do i hate kubernetes");


export default function (client: Client) {
    program.command("deploy")
        .description("Deploys whatever the fuck you want, easily")
        .option(
            "-f, --file <string>",
            "Specify a project file to deploy",
            "hive.yml",
        )
        .action((options) => client.deploy(options));
    program.parse();
}