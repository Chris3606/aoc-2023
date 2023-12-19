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

// ReSharper disable NotAccessedPositionalProperty.Global
public readonly record struct SpringRecordState(int RecordIdx, char RecordValue, int GroupIdx);
// ReSharper restore NotAccessedPositionalProperty.Global

public sealed class Day12 : BaseDay
{
    private readonly SpringRecord[] _records;

    public Day12()
    {
        _records = File.ReadLines(InputFilePath).Select(SpringRecord.Parse).ToArray();
    }

    public override ValueTask<string> Solve_1()
        => new(_records.Sum(record => CountRecord(record.Springs, 0, record.Groups, 0)).ToString());

    public override ValueTask<string> Solve_2()
        => new(_records.Select(i => Unfold(i, 5)).Sum(record => CountRecord(record.Springs, 0, record.Groups, 0)).ToString());

    private long CountRecord(Span<char> record, int recordIdx, ReadOnlySpan<int> groups, int groupIdx, Dictionary<SpringRecordState, long> memoTable = null)
    {
        if (record.Length == 0)
            return groups.Length == 0 ? 1 : 0;

        memoTable ??= new Dictionary<SpringRecordState, long>();
        
        var state = new SpringRecordState(recordIdx, record[0], groupIdx);
        if (memoTable.TryGetValue(state, out long memoized))
            return memoized;
        
        
        long permutations = 0;
        switch (record[0])
        {
            case '.':
                permutations = CountRecord(record[1..], recordIdx + 1, groups, groupIdx, memoTable);
                break;
            case '?':
                record[0] = '#';
                permutations += CountRecord(record, recordIdx, groups, groupIdx, memoTable);
                
                record[0] = '.';
                permutations += CountRecord(record[1..], recordIdx + 1, groups, groupIdx, memoTable);

                record[0] = '?';
                break;
            case '#':
                if (groups.Length == 0)
                    break;

                if (record.Length < groups[0])
                    break;

                if (SpanExtensions.Count(record[..groups[0]], '.') > 0)
                    break;

                if (record.Length == groups[0])
                {
                    permutations = groups.Length == 1 ? 1 : 0;
                    break;
                }


                permutations = record[groups[0]] == '#' ? 0 : CountRecord(record[(groups[0] + 1)..], recordIdx + groups[0] + 1, groups[1..], groupIdx + 1, memoTable);
                break;
            
            default:
                throw new ArgumentException("Unsupported spring record", nameof(record));
        }
        
        memoTable.Add(state, permutations);

        return permutations;
    }

    private static SpringRecord Unfold(SpringRecord record, int times)
    {
        IEnumerable<char> springs = record.Springs;
        IEnumerable<int> groups = record.Groups;
        for (int i = 1; i < times; i++)
        {
            springs = springs.Concat('?'.Yield()).Concat(record.Springs);
            groups = groups.Concat(record.Groups);
        }

        return new(springs.ToArray(), groups.ToArray());
    }
}
