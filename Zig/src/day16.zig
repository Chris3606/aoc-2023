const std = @import("std");

const util = @import("util.zig");
const gpa = util.gpa;

const data = @embedFile("data/day16.txt");

fn part1() !i32 {
    return -1;
}

fn part2() !i32 {
    return -1;
}

pub fn main() !void {
    defer {
        const deinit_status = util.gpa_impl.deinit();
        if (deinit_status == .leak) @panic("Leaked memory!");
    }

    std.log.info("Part 1: {}", .{try part1()});
    std.log.info("Part 2: {}", .{try part2()});
}

// Generated from template/template.zig.
// Run `zig build generate` to update.
// Only unmodified days will be updated.
