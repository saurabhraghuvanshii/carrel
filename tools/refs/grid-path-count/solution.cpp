#include <vector>
using namespace std;

int countPaths(vector<vector<int>>& grid) {
    int cols = grid[0].size();
    vector<int> ways(cols, 0);
    ways[0] = grid[0][0] == 0;
    for (auto& row : grid)
        for (int c = 0; c < cols; c++) {
            if (row[c] == 1) ways[c] = 0;
            else if (c > 0) ways[c] += ways[c - 1];
        }
    return ways[cols - 1];
}
