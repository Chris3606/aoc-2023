const std = @import("std");
const Allocator = std.mem.Allocator;
const List = std.ArrayList;
const Map = std.AutoHashMap;
const StrMap = std.StringHashMap;
const BitSet = std.DynamicBitSet;
const Str = []const u8;

pub var gpa_impl = std.heap.GeneralPurposeAllocator(.{}){};
pub const gpa = gpa_impl.allocator();

// Add utility functions here

// Parsing errors
pub const ParseError = error {
    InvalidData
};

// Iterates forward over a slice
fn ForwardIterator(comptime T: type) type {
    const Pointer = blk: {
        switch (@typeInfo(T)) {
            .Pointer => |ptr_info| switch (ptr_info.size) {
                .One => switch (@typeInfo(ptr_info.child)) {
                    .Array => |array_info| {
                        var new_ptr_info = ptr_info;
                        new_ptr_info.size = .Many;
                        new_ptr_info.child = array_info.child;
                        new_ptr_info.sentinel = array_info.sentinel;
                        break :blk @Type(.{ .Pointer = new_ptr_info });
                    },
                    else => {},
                },
                .Slice => {
                    var new_ptr_info = ptr_info;
                    new_ptr_info.size = .Many;
                    break :blk @Type(.{ .Pointer = new_ptr_info });
                },
                else => {},
            },
            else => {},
        }
        @compileError("expected slice or pointer to array, found '" ++ @typeName(T) ++ "'");
    };
    const Element = std.meta.Elem(Pointer);
    const ElementPointer = @Type(.{ .Pointer = ptr: {
        var ptr = @typeInfo(Pointer).Pointer;
        ptr.size = .One;
        ptr.child = Element;
        ptr.sentinel = null;
        break :ptr ptr;
    } });
    return struct {
        ptr: Pointer,
        len: usize,
        index: usize,
        pub fn next(self: *@This()) ?Element {
            if (self.index == self.len) return null;
            const cur = self.index;
            self.index += 1;

            return self.ptr[cur];
        }
        pub fn nextPtr(self: *@This()) ?ElementPointer {
            if (self.index == self.len) return null;
            const cur = self.index;
            self.index += 1;

            return &self.ptr[cur];
        }
    };
}

/// Iterates over a slice.
pub fn sliceIterator(slice: anytype) ForwardIterator(@TypeOf(slice)) {
    return .{ .ptr = slice.ptr, .index = 0, .len = slice.len };
}

// Useful stdlib functions
const tokenizeAny = std.mem.tokenizeAny;
const tokenizeSeq = std.mem.tokenizeSequence;
const tokenizeSca = std.mem.tokenizeScalar;
const splitAny = std.mem.splitAny;
const splitSeq = std.mem.splitSequence;
const splitSca = std.mem.splitScalar;
const indexOf = std.mem.indexOfScalar;
const indexOfAny = std.mem.indexOfAny;
const indexOfStr = std.mem.indexOfPosLinear;
const lastIndexOf = std.mem.lastIndexOfScalar;
const lastIndexOfAny = std.mem.lastIndexOfAny;
const lastIndexOfStr = std.mem.lastIndexOfLinear;
const trim = std.mem.trim;
const sliceMin = std.mem.min;
const sliceMax = std.mem.max;

const parseInt = std.fmt.parseInt;
const parseFloat = std.fmt.parseFloat;

const print = std.debug.print;
const assert = std.debug.assert;

const sort = std.sort.block;
const asc = std.sort.asc;
const desc = std.sort.desc;
