#include <string>
#include <utility>
#include <vector>
using namespace std;

int countIslands(vector<string>& grid) {
    int rows = grid.size(), cols = grid[0].size(), count = 0;
    for (int r = 0; r < rows; r++)
        for (int c = 0; c < cols; c++) {
            if (grid[r][c] != '1') continue;
            count++;
            vector<pair<int, int>> stack = {{r, c}};
            grid[r][c] = '0';
            while (!stack.empty()) {
                auto [y, x] = stack.back();
                stack.pop_back();
                int dy[] = {0, 1, 0, -1}, dx[] = {1, 0, -1, 0};
                for (int d = 0; d < 4; d++) {
                    int ny = y + dy[d], nx = x + dx[d];
                    if (ny >= 0 && ny < rows && nx >= 0 && nx < cols && grid[ny][nx] == '1') {
                        grid[ny][nx] = '0';
                        stack.push_back({ny, nx});
                    }
                }
            }
        }
    return count;
}
