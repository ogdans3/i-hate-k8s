const std = @import("std");
const cli = @import("zig-cli");

// Define a configuration structure with default values.
var config = struct {
    host: []const u8 = "localhost",
    port: u16 = undefined,
}{};

pub fn command_line_interface() !void {
    var r = try cli.AppRunner.init(std.heap.page_allocator);

    // Create an App with a command named "short" that takes host and port options.
    const app = cli.App{
        .command = cli.Command{
            .name = "short",
            .options = &.{
                // Define an Option for the "host" command-line argument.
                .{
                    .long_name = "host",
                    .help = "host to listen on",
                    .value_ref = r.mkRef(&config.host),
                },

                // Define an Option for the "port" command-line argument.
                .{
                    .long_name = "port",
                    .help = "port to bind to",
                    .required = true,
                    .value_ref = r.mkRef(&config.port),
                },

            },
            .target = cli.CommandTarget{
                .action = cli.CommandAction{ .exec = run_server },
            },
        },
    };
    return r.run(&app);
}

// Action function to execute when the "short" command is invoked.
fn run_server() !void {
    // Log a debug message indicating the server is listening on the specified host and port.
    std.log.debug("server is listening on {s}:{d}", .{ config.host, config.port });
}