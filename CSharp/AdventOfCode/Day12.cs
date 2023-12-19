using CommunityToolkit.HighPerformance;

namespace AdventOfCode;

public record SpringRecord(char[] Springs, int[] Groups)
{
    public static SpringRecord Parse(string line)
    {
        var parts = line.Split(" ");
        return new(parts[0].ToArray(), parts[1].Split(",").Select(int.Parse).ToArray());
    }
}

public sealed class Day12 : BaseDay
{
    private readonly SpringRecord[] _records;

    public Day12()
    {
        _records = File.ReadLines(InputFilePath).Select(SpringRecord.Parse).ToArray();
    }

    public override ValueTask<string> Solve_1()
        => new(_records.Sum(record => CountRecord(record.Springs, record.Groups)).ToString());

    public override ValueTask<string> Solve_2() => throw new NotImplementedException();

    private int CountRecord(Span<char> record, ReadOnlySpan<int> groups)
    {
        if (record.Length == 0)
            return groups.Length == 0 ? 1 : 0;
        
        switch (record[0])
        {
            case '.':
                return CountRecord(record[1..], groups);
            case '?':
                int permutations = 0;
                record[0] = '#';
                permutations += CountRecord(record, groups);
                
                record[0] = '.';
                permutations += CountRecord(record[1..], groups);

                record[0] = '?';
                return permutations;
            case '#':
                if (groups.Length == 0)
                    return 0;
                
                if (record.Length < groups[0])
                    return 0;
                
                if (SpanExtensions.Count(record[..groups[0]], '.') > 0)
                    return 0;

                if (record.Length == groups[0])
                    return groups.Length == 1 ? 1 : 0;
                
                
                return record[groups[0]] == '#' ? 0 : CountRecord(record[(groups[0] + 1)..], groups[1..]);

            default:
                throw new ArgumentException("Unsupported spring record", nameof(record));
        }
    }
}
