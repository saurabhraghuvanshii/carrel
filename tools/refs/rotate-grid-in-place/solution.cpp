#include <algorithm>
#include <vector>
using namespace std;

void rotateGrid(vector<vector<int>>& grid) {
    int n = grid.size();
    for (int r = 0; r < n; r++)
        for (int c = r + 1; c < n; c++) swap(grid[r][c], grid[c][r]);
    for (auto& row : grid) reverse(row.begin(), row.end());
}
