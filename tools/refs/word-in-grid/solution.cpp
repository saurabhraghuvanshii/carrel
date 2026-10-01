#include <string>
#include <vector>
using namespace std;

static bool walk(vector<string>& g, const string& word, int r, int c, size_t k) {
    if (r < 0 || r >= (int)g.size() || c < 0 || c >= (int)g[0].size() || g[r][c] != word[k]) return false;
    if (k + 1 == word.size()) return true;
    char keep = g[r][c];
    g[r][c] = '#';
    bool found = walk(g, word, r + 1, c, k + 1) || walk(g, word, r - 1, c, k + 1) ||
                 walk(g, word, r, c + 1, k + 1) || walk(g, word, r, c - 1, k + 1);
    g[r][c] = keep;
    return found;
}

bool hasWord(vector<string>& grid, string& word) {
    for (int r = 0; r < (int)grid.size(); r++)
        for (int c = 0; c < (int)grid[0].size(); c++)
            if (walk(grid, word, r, c, 0)) return true;
    return false;
}
