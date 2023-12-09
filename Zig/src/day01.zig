const std = @import("std");

const util = @import("util.zig");
const gpa = util.gpa;

const data = @embedFile("data/day01.txt");

const words = [_][]const u8{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"};

fn findFirstNumber(str: []const u8, include_words: bool) !u8 {

    for (0.., str) |idx, val| {
        if (include_words) {
            for (0.., words) |wordVal, word| {
                if (std.mem.endsWith(u8, str[0..idx], word)) {
                    return @as(u8, @intCast(wordVal)) + '0';
                }
            }
        }

        if (std.ascii.isDigit(val)) {
            return val;
        }
    }

    return error.InvalidData;
}

fn findLastNumber(str: []const u8, include_words: bool) !u8 {

    var idx = str.len;
    while (idx > 0) {
        idx-= 1;

        if (include_words) {
            for (0.., words) |wordVal, word| {
                if (std.mem.startsWith(u8, str[idx..], word)) {
                    return @as(u8, @intCast(wordVal)) + '0';
                }
            }
        }

        if (std.ascii.isDigit(str[idx])) {
            return str[idx];
        }
    }

    return error.InvalidData;
}

fn part1() !i32 {
    var it = std.mem.tokenizeScalar(u8, data, '\n');

    var sum: i32 = 0;
    while (it.next()) |line| {
        var num: [2]u8 = undefined;

        num[0] = try findFirstNumber(line, false);
        num[1] = try findLastNumber(line, true);

        sum += try std.fmt.parseInt(i32, &num, 10);
    }

    return sum;
}

fn part2() !i32 {
    var it = std.mem.tokenizeScalar(u8, data, '\n');

    var sum: i32 = 0;
    while (it.next()) |line| {
        var num: [2]u8 = undefined;
        num[1] = '1';

        num[0] = try findFirstNumber(line, true);
        num[1] = try findLastNumber(line, true);

        sum += try std.fmt.parseInt(i32, &num, 10);
    }

    return sum;
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
