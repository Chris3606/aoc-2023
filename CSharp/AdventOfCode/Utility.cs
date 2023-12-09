namespace AdventOfCode;

public static class Utility
{
    public static Dictionary<T, int> BuildHistogram<T>(this IEnumerable<T> self)
    {
        var dict = new Dictionary<T, int>();
        foreach (var item in self)
        {
            if (!dict.TryAdd(item, 1))
                dict[item]++;
        }

        return dict;
    }
}