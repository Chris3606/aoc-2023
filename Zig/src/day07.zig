const std = @import("std");

const util = @import("util.zig");
const gpa = util.gpa;

const data = @embedFile("data/day07.txt");

pub fn main() !void {
    defer {
        const deinit_status = util.gpa_impl.deinit();
        if (deinit_status == .leak) @panic("Leaked memory!");
    }
}

// Generated from template/template.zig.
// Run `zig build generate` to update.
// Only unmodified days will be updated.
