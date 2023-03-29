import { ShardingManager, fetchRecommendedShardCount } from 'discord.js';
import log from './utils/log';

const shardManager = new ShardingManager('./sdc.js', {
	token: process.env.TOKEN,
	execArgv: ['--unhandled-rejections=warn', '--experimental-specifier-resolution=node'],
});

shardManager.on('shardCreate', shard => log(`Shard spawned!`, shard.id));

(async () => {
	const amount = await fetchRecommendedShardCount(process.env.TOKEN, {
		guildsPerShard: 1000,
	});
	log(`Shards count: ${amount}`, 'null');

	shardManager.shardList = [...Array(amount).keys()];
	shardManager.totalShards = amount;

	// Spawn the shards
	const promises = [];
	for (let i = 0; i < amount; i++) promises.push(shardManager.createShard(i).spawn(5 * 60e3));
	const shards = await Promise.all(promises);

	for (const shard of shards) shard.send('startPresence');
})();
