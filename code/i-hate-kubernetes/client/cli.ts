import { Command } from "@commander-js/extra-typings";
const program = new Command();
export default program;

program
    .name("I hate kubernetes")
    .description("Boy do i hate kubernetes");

program.command("deploy")
    .description("Deploys whatever the fuck you want, easily")
    .option(
        "-f, --file <string>",
        "Specify a project file to deploy",
        "hive.yml",
    )
    .action((options) => deploy(options));

function deploy(options: { file: string }) {
    const { file } = options;
    console.log(file);
}
