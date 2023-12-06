using SadRogue.Primitives;
using SadRogue.Primitives.GridViews;

namespace AdventOfCode;

public class Number(Point position)
{
    public int Value { get; private set; }
    public int NumDigits { get; private set; }
    
    public Point Position { get; private set; } = position;

    public void AddDigit(int digit)
    {
        Value *= 10;
        Value += digit;
        NumDigits++;
    }
}

public sealed class Day03 : BaseDay
{
    private readonly ArrayView<char> _grid;
    private readonly Dictionary<Point, Number> _numbers;

    public static Dictionary<Point, Number> GetNumbers(IGridView<char> grid)
    {
        var nums = new Dictionary<Point, Number>();

        Number curNum = null;
        for (int y = 0; y < grid.Height; y++)
        {
            for (int x = 0; x < grid.Width; x++)
            {
                var p = new Point(x, y);
                
                var c = grid[x, y];
                if (c is >= '0' and <= '9')
                {
                    curNum ??= new Number(p);
                    curNum.AddDigit(c - '0');
                    nums[p] = curNum;
                }
                else
                    curNum = null;
            }

            curNum = null;
        }

        return nums;
    }

    public Day03()
    {
        _grid = ParseUtils.CharGridFrom(File.ReadAllText(InputFilePath));
        _numbers = GetNumbers(_grid);
    }

    public override ValueTask<string> Solve_1()
    {
        var partNumbers = _numbers.Where(pair
            => AdjacencyRule.EightWay.Neighbors(pair.Key).Any(n
                => !_numbers.ContainsKey(n) && _grid.Contains(n) && _grid[n] is < '0' or > '9' && _grid[n] != '.'))
            .Select(pair => pair.Value)
            .ToHashSet();

        return new(partNumbers.Sum(i => i.Value).ToString());
    }

    public override ValueTask<string> Solve_2()
    {
        var answer = _grid.Positions()
            .Where(p => _grid[p] == '*')
            .Select(point =>
            {
                var hs = new HashSet<Number>();
                foreach (var n in AdjacencyRule.EightWay.Neighbors(point)
                             .Where(p => _numbers.ContainsKey(p)))
                    hs.Add(_numbers[n]);

                return hs;
            })
            .Where(h => h.Count == 2)
            .Select(i => i.Aggregate(1, (a, num) => a * num.Value))
            .Sum();
        
        return new ValueTask<string>(answer.ToString());
    }
}
